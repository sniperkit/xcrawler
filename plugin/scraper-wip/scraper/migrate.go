package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	// "log"
	"net"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/Machiel/slugify"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
	"github.com/roscopecoltran/e3ch"
	"github.com/roscopecoltran/mxj"
	// "github.com/ksinica/flatstruct"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/shurcooL/go-goon"
	// "github.com/kr/pretty"
)

var (
	slugifier = slugify.New(slugify.Configuration{
		ReplaceCharacter: '_',
	})
)

func MigrateTables(db *gorm.DB, isTruncate bool, tables ...interface{}) {
	for _, table := range tables {
		if isTruncate {
			if err := db.DropTableIfExists(table).Error; err != nil {
				fmt.Println("table creation error, error msg: ", err)
			}
		}
		db.AutoMigrate(table)
	}
}

// FindOrCreateTagByName finds a tag by name, creating if it doesn't exist
func FindOrCreateProviderByName(db *gorm.DB, name string) (Provider, bool, error) {
	if name == "" {
		return Provider{}, false, errors.New("WARNING !!! No provider name provided")
	}
	var provider Provider
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&provider).RecordNotFound() {
		provider.Name = name
		err := db.Create(&provider).Error
		return provider, true, err
	}
	return provider, false, nil
}

func FindOrCreateGroupByName(db *gorm.DB, name string) (*Group, bool, error) {
	if name == "" {
		return nil, false, errors.New("WARNING !!! No provider name provided")
	}
	var group Group
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&group).RecordNotFound() {
		group.Name = name
		err := db.Create(&group).Error
		return &group, true, err
	}
	return &group, false, nil
}

func MigrateEndpoints(db *gorm.DB, c Config, e3ch *client.EtcdHRCHYClient) error {
	for _, e := range c.Routes {
		selectionBlocks, err := convertSelectorsConfig(e.BlocksJSON, c.Debug)
		if err != nil {
			return err
		}
		headers, err := convertHeadersConfig(e.HeadersJSON, c.Debug)
		if err != nil {
			return err
		}
		endpointTemplateURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(e.BaseURL, "/"), strings.TrimPrefix(e.PatternURL, "/"))
		slugURL := slugifier.Slugify(endpointTemplateURL)
		var groups []*Group
		group, _, err := FindOrCreateGroupByName(db, "Web")
		if err != nil {
			if c.Debug {
				fmt.Println("Could not upsert the group for the current endpoint. error: ", err)
			}
			//return err
		}
		groups = append(groups, group)
		providerDataURL, err := url.Parse(e.BaseURL)
		if err != nil {
			if c.Debug {
				fmt.Println("Could not parse/extract the endpoint url parts. error: ", err)
			}
			//return err
		}
		providerHost, providerPort, _ := net.SplitHostPort(providerDataURL.Host)
		//if err != nil {
		//	fmt.Println("Could not split host and port for the current endpoint base url. error: ", err)
		// return err
		//}
		providerDomain := domainutil.Domain(providerDataURL.Host)
		provider, _, err := FindOrCreateProviderByName(db, providerDomain)
		if err != nil {
			if c.Debug {
				fmt.Println("Could not upsert the current provider in the registry. error: ", err)
			}
			// return err
		}

		if e3ch != nil {
			etcdHeadersIntercept := make(map[string]string, 0)
			for _, v := range e.HeadersIntercept {
				if v != "" {
					if c.Debug {
						fmt.Println("new etcdHeadersIntercept value: ")
					}
					parts := strings.Split(v, ":")
					if len(parts) > 1 {
						key := parts[0]
						val := parts[1]
						etcdHeadersIntercept[key] = val
					}
				}
			}
			etcdHeaders := make(map[string]string, 0)
			for k, v := range e.HeadersJSON {
				if c.Debug {
					fmt.Println("new etcdHeaders value: ")
				}
				etcdHeaders[k] = v
			}
			etcdBlocks := make(map[string]map[string]string, 0)
			for k, v := range e.BlocksJSON {
				etcdBlocks[k] = make(map[string]string, 0)
				etcdBlocks[k]["items"] = v.Items
				for kd, vd := range v.Details {
					var ext []string
					for _, vdd := range vd {
						ext = append(ext, vdd.val)
					}
					etcdBlocks[k][kd] = strings.Join(ext, ";")
				}
				if c.Debug {
					fmt.Println("new etcdBlocks value: ")
					pp.Println(v)
				}
			}

			var ExtractorTypes = []string{"links", "meta", "opengraph"}
			etcdExtract := make(map[string]bool, len(ExtractorTypes))
			for _, v := range ExtractorTypes {
				etcdExtract[v] = false
			}

			var etcdGroups string
			var grp []string
			for _, v := range e.Groups {
				grp = append(grp, v.Name)
			}
			etcdGroups = strings.Join(grp, ",")
			etcdRoute := EtcdRoute{
				Loaded:           true,
				Disabled:         false,
				Source:           provider.Name,
				Route:            e.Route,
				Method:           strings.ToUpper(e.Method),
				BaseURL:          e.BaseURL,
				PatternURL:       e.PatternURL,
				TestURL:          endpointTemplateURL,
				Selector:         e.Selector,
				HeadersIntercept: etcdHeadersIntercept,
				Headers:          etcdHeaders,
				Comment:          "-",
				Blocks:           etcdBlocks,
				Groups:           etcdGroups,
				Extract:          etcdExtract,
				StrictMode:       false,
				Debug:            false,
			}

			b, err := json.Marshal(etcdRoute)
			if err != nil {
				fmt.Println(err)
				// return
			} else {
				// fmt.Println("new etcdRoute: ")
				// pp.Print(etcdRoute)
				mv, err := mxj.NewMapJson(b)
				if err != nil {
					fmt.Println("err:", err)
					// return
				} else {
					mxj.LeafSeparator = "/"
					// mxj.LeafUseDotNotation()
					p := mv.LeafPaths()
					var tree []string
					for _, v := range p {
						p := strings.Split(fmt.Sprintf("/endpoint/%s/config/scraper/%s", e.Route, v), "/")
						j := len(p) - 1
						for j > 1 {
							tree = append(tree, strings.Join(p[:j], "/"))
							// fmt.Printf("sub keys: %s \n", strings.Join(p[:j], "/"))
							j--
						}
					}
					// dedup(tree)
					// tree = removeDuplicates(tree)
					RemoveDuplicates(&tree)
					sort.Strings(tree)
					if c.Debug {
						fmt.Println("LeafPaths for ETCD: ")
						pp.Println(tree)
					}
					for _, d := range tree {
						if c.Debug {
							fmt.Println("create dir: ", d)
						}
						if d != "" {
							err := e3ch.CreateDir(d)
							if err != nil {
								if c.Debug {
									fmt.Printf("Could not delete dir '%s', error: %s \n", d, err)
								}
							}
						}
					}
					l := mv.LeafNodes()
					if c.Debug {
						fmt.Println("LeafNodes: ")
					}
					for _, v := range l {
						key := fmt.Sprintf("/endpoint/%s/config/scraper/%s", e.Route, v.Path)
						val := fmt.Sprintf("%v", v.Value)
						if c.Debug {
							fmt.Printf("path: %s, value: %s \n", key, val)
						}
						var exists bool
						if val != "" {
							if c.Debug {
								fmt.Println("create key: ", key)
							}
							err := e3ch.Create(key, val)
							if err != nil {
								if c.Debug {
									fmt.Printf("Could not create key key='%s', value='%s', error: %s \n", key, val, err)
								}
								exists = true
							}
							if exists {
								if c.Debug {
									fmt.Println("put/update key: ", key)
								}
								err := e3ch.Put(key, val)
								if err != nil {
									if c.Debug {
										fmt.Printf("Could not put key='%s', value='%s', error: %s \n", key, val, err)
									}
								}
							}
						}
					}
				}
			}
		}

		endpoint := Endpoint{
			Disabled:     false,
			Route:        e.Route,
			Method:       strings.ToUpper(e.Method),
			BaseURL:      e.BaseURL,
			PatternURL:   e.PatternURL,
			Selector:     e.Selector,
			Slug:         slugURL,
			Headers:      headers,
			Blocks:       selectionBlocks,
			ExtractPaths: e.ExtractPaths,
			Debug:        e.Debug,
			StrictMode:   e.StrictMode,
		}
		endpoint.Groups = groups
		if providerHost != "" {
			endpoint.Host = providerHost
		} else {
			endpoint.Host = providerDataURL.Host
		}
		endpoint.Domain = providerDomain
		if providerPort != "" {
			providerPortInt, err := strconv.Atoi(providerPort)
			if err != nil {
				if c.Debug {
					fmt.Println("WARNING ! Missing the port number for this endpoint base url. error: ", err)
				}
			}
			endpoint.Port = providerPortInt
		} else {
			// Move to a seperate method
			switch providerDataURL.Scheme {
			case "wss":
			case "https":
				endpoint.Port = 443
			case "ws":
			case "http":
				endpoint.Port = 80
			case "rpc":
				endpoint.Port = 445
			default:
				if c.Debug {
					fmt.Println("WARNING ! invalid base url scheme for the current endpoint.")
				}
			}
		}
		endpoint.Description = e.Description
		endpoint.Provider = provider
		for _, b := range selectionBlocks {
			if ok := db.NewRecord(b); ok {
				if err := db.Create(&b).Error; err != nil {
					if c.Debug {
						fmt.Println("error: ", err)
					}
					return err
				}
			}
		}
		if ok := db.NewRecord(endpoint); ok {
			if err := db.Create(&endpoint).Error; err != nil {
				if c.Debug {
					fmt.Println("error: ", err)
				}
				return err
			}
		}
		e.ready = true
		e.hash, err = e.getHash("sha1")
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func convertProviderConfig(name string, debug bool) *Provider {
	provider := &Provider{}
	if name != "" {
		provider.Name = name
	} else {
		return nil
	}
	if debug {
		fmt.Printf("\nConverting provider name: '%s' \n", name)
	}
	return provider
}

func convertSelectorsConfig(selectors map[string]*SelectorConfig, debug bool) ([]*SelectorConfig, error) {
	var blocks []*SelectorConfig
	for k, v := range selectors {
		targets, err := convertDetailsConfig(v.Details, debug)
		if err != nil {
			return nil, err
		}
		selection := &SelectorConfig{
			Collection: k,
			Debug:      v.Debug,
			Required:   v.Required,
			Items:      v.Items,
			Matchers:   targets,
			StrictMode: v.StrictMode,
		}
		blocks = append(blocks, selection)
	}
	return blocks, nil
}

func convertDetailsConfig(tgts map[string]Extractors, debug bool) ([]*MatcherConfig, error) {
	var targets []*MatcherConfig
	for k, t := range tgts {
		var matchers []Matcher
		for _, e := range t {
			matchers = append(matchers, Matcher{Expression: e.val})
		}
		target := &MatcherConfig{
			Target:  k,
			Selects: matchers,
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func convertHeadersConfig(headers map[string]string, debug bool) ([]*HeaderConfig, error) {
	var hdrs []*HeaderConfig
	for k, v := range headers {
		header := &HeaderConfig{
			Key:   k,
			Value: v,
		}
		if debug {
			fmt.Printf("\nConverting header config: %s:%s \n", k, v)
		}
		hdrs = append(hdrs, header)
	}
	return hdrs, nil
}

func createGroups(db *gorm.DB) {
	for _, g := range Seeds.Groups {
		group := Group{}
		group.Name = g.Name
		if err := db.Create(&group).Error; err != nil {
			fmt.Printf("create group (%v) failure, got err %v\n", group, err)
		}
	}
}

func createTopics(db *gorm.DB) {
	for _, t := range Seeds.Topics {
		topic := Topic{}
		topic.Name = t.Name
		topic.Code = strings.ToLower(t.Name)
		if err := db.Create(&topic).Error; err != nil {
			fmt.Printf("create topic (%v) failure, got err %v\n", topic, err)
		}
	}
}

/*
func convertSelectorsToEtcd(selectors map[string]SelectorConfig, debug bool) {
	var blocks []*SelectorConfig
	for k, v := range selectors {
		targets, err := convertDetailsToEtcd(v.Details, debug)
		if err != nil {
			return nil, err
		}
		selection := &SelectorConfig{
			Collection: k,
			Debug:      v.Debug,
			Required:   v.Required,
			Items:      v.Items,
			Matchers:   targets,
			StrictMode: v.StrictMode,
		}
		blocks = append(blocks, selection)
	}
	return blocks, nil
}

func convertDetailsToEtcd(tgts map[string]Extractors, debug bool) ([]*MatcherConfig, error) {
	var targets []*MatcherConfig
	for k, t := range tgts {
		var matchers []Matcher
		for _, e := range t {
			matchers = append(matchers, Matcher{Expression: e.val})
		}
		target := &MatcherConfig{
			Target:  k,
			Selects: matchers,
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func convertHeadersConfig(headers map[string]string, debug bool) {

}
*/
