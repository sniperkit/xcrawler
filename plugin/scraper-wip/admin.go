package main

import (
	"github.com/qor/action_bar"
	"github.com/roscopecoltran/scraper/scraper"

	admin_help "github.com/qor/help"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/roscopecoltran/admin"
)

var (
	ActionBar *action_bar.ActionBar
)

func initDashboard() {

	// Initalize
	// AdminUI = admin.New(&qor.Config{DB: db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)})
	AdminUI = admin.New(&qor.Config{DB: DB})

	// Meta info
	AdminUI.SetSiteName("Sniperkit-Scraper Config")

	// Auth
	// AdminUI.SetAuth(auth.AdminAuth{})

	// Assets FileSystem
	// AdminUI.SetAssetFS(bindatafs.AssetFS)

	// Menu(s)
	AdminUI.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin"}) // // Add Dashboard

	// Categories (Scrapers, Providers)
	// topic := AdminUI.AddResource(&scraper.Topic{}) //, &admin.Config{Menu: []string{"Source Management"}})
	// topic.Meta(&admin.Meta{Name: "Topics", Type: "select_many"})

	// category := Admin.AddResource(&models.Category{}, &admin.Config{Menu: []string{"Product Management"}, Priority: -3})
	// category.Meta(&admin.Meta{Name: "Categories", Type: "select_many"})

	// m := qor_admin_auth.DefaultNew(DB, &models.User{})
	// m.AddResourceForPage("user", &models.User{}, "User Manage")

	// Activity
	queries := AdminUI.AddResource(&scraper.Query{}, &admin.Config{Menu: []string{"Activity"}})
	//query := queries.Meta(&admin.Meta{Name: "Keywords"}).Resource
	queries.NewAttrs(&admin.Section{
		Rows: [][]string{{"InputQuery", "Blocked"}},
	})
	queries.EditAttrs(&admin.Section{
		Rows: [][]string{{"InputQuery", "Blocked"}},
	})

	// Groups of Scrapers
	group := AdminUI.AddResource(&scraper.Group{}, &admin.Config{Menu: []string{"Classify Data"}})
	// group.Meta(&admin.Meta{Name: "Groups", Type: "select_many"})

	// Add Asset Manager, for rich editor
	assetManager := AdminUI.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})

	// Add Help
	Help := AdminUI.NewResource(&admin_help.QorHelpEntry{}, &admin.Config{Menu: []string{"Help"}})
	Help.GetMeta("Body").Config = &admin.RichEditorConfig{AssetManager: assetManager}

	details := AdminUI.AddResource(&scraper.MatcherConfig{}, &admin.Config{Invisible: true})
	details.Meta(&admin.Meta{Name: "Target", Config: &admin.SelectOneConfig{Collection: scraper.TargetTypes, AllowBlank: false}})

	AdminUI.AddResource(&scraper.SelectorConfig{}, &admin.Config{Invisible: true})
	AdminUI.AddResource(&scraper.HeaderConfig{}, &admin.Config{Invisible: true})

	// connection :=
	connection := AdminUI.AddResource(&scraper.Connection{}, &admin.Config{Menu: []string{"Activity"}})
	connection.IndexAttrs("ID", "Provider", "URL", "Response.Code")
	// connection.Meta(&admin.Meta{Name: "Name", Type: "Readonly"})

	AdminUI.AddResource(&scraper.Request{}, &admin.Config{Invisible: true})
	// request.Meta(&admin.Meta{Name: "Body", Type: "text"})
	// request.IndexAttrs("ID", "VersionName", "ScheduledStartAt", "ScheduledEndAt", "Author", "Title")

	AdminUI.AddResource(&scraper.Response{}, &admin.Config{Invisible: true})

	// Endpoints
	endpoint := AdminUI.AddResource(&scraper.Endpoint{}, &admin.Config{Menu: []string{"Web Scrapers"}})
	endpoint.Meta(&admin.Meta{Name: "Selector", Config: &admin.SelectOneConfig{Collection: scraper.SelectorEngines, AllowBlank: false}})
	endpoint.Meta(&admin.Meta{Name: "Method", Config: &admin.SelectOneConfig{Collection: scraper.MethodTypes, AllowBlank: false}})
	endpoint.Meta(&admin.Meta{Name: "Groups", Config: &admin.SelectManyConfig{SelectMode: "bottom_sheet"}})

	endpoint.IndexAttrs("Name", "Disabled", "Provider", "Route", "Method")
	endpoint.SearchAttrs("Name", "Disabled", "Provider", "Route", "Method")

	endpoint.Meta(&admin.Meta{Name: "Description", Config: &admin.RichEditorConfig{AssetManager: assetManager, Plugins: []admin.RedactorPlugin{
		{Name: "medialibrary", Source: "/admin/assets/javascripts/qor_redactor_medialibrary.js"},
		{Name: "table", Source: "/vendors/redactor_table.js"},
	},
		Settings: map[string]interface{}{
			"medialibraryUrl": "/admin/product_images",
		},
	}})

	endpoint.Filter(&admin.Filter{
		Name:   "Groups",
		Config: &admin.SelectOneConfig{RemoteDataResource: group},
	})

	// Providers
	provider := AdminUI.AddResource(&scraper.Provider{}, &admin.Config{Menu: []string{"Classify Data"}})
	providerWebRank := provider.Meta(&admin.Meta{Name: "Ranks"}).Resource
	providerWebRank.ShowAttrs("Engine", "Score")

	endpoint.Filter(&admin.Filter{
		Name:   "Providers",
		Config: &admin.SelectOneConfig{RemoteDataResource: provider},
	})

	// product.SearchAttrs("Name", "Code", "Category.Name", "Brand.Name")

	headersEndpoint := endpoint.Meta(&admin.Meta{Name: "Headers"}).Resource
	headersEndpoint.NewAttrs(&admin.Section{
		Rows: [][]string{{"Key", "Value"}},
	})
	headersEndpoint.EditAttrs(&admin.Section{
		Rows: [][]string{{"Key", "Value"}},
	})

	blocksEndpoint := endpoint.Meta(&admin.Meta{Name: "Blocks"}).Resource
	blocksEndpoint.EditAttrs("Name", "Disabled", "Items", "Required", "Description", "Matchers", "StrictMode", "Debug")

	blocksEndpoint.Meta(&admin.Meta{Name: "Description", Config: &admin.RichEditorConfig{AssetManager: assetManager, Plugins: []admin.RedactorPlugin{
		{Name: "medialibrary", Source: "/admin/assets/javascripts/qor_redactor_medialibrary.js"},
		{Name: "table", Source: "/vendors/redactor_table.js"},
	},
		Settings: map[string]interface{}{
			"medialibraryUrl": "/admin/product_images",
		},
	}})

	detailsEndpoint := blocksEndpoint.Meta(&admin.Meta{Name: "Matchers"}).Resource
	detailsEndpoint.Meta(&admin.Meta{Name: "Target", Config: &admin.SelectOneConfig{Collection: scraper.TargetTypes, AllowBlank: false}})

	// Add ProductImage as Media Libraray
	ScreenshotsResource := AdminUI.AddResource(&scraper.Screenshot{}, &admin.Config{Menu: []string{"Activity"}, Priority: -1})
	ScreenshotsResource.Filter(&admin.Filter{
		Name:       "SelectedType",
		Label:      "Media Type",
		Operations: []string{"contains"},
		Config:     &admin.SelectOneConfig{Collection: [][]string{{"video", "Video"}, {"image", "Image"}, {"file", "File"}, {"video_link", "Video Link"}}},
	})
	ScreenshotsResource.IndexAttrs("File", "Title")

	// endpoint.ShowAttrs("Disabled", "Debug", "Name", "Route", "Method", "ExampleURL", "Selector", "BaseURL", "PatternURL", "Headers", "Blocks", "Extract", "StrictMode")
	/*
		endpoint.NewAttrs(
			&admin.Section{
				Title: "Status",
				Rows: [][]string{
					{"Disabled", "Debug"},
				},
			},
			&admin.Section{
				Title: "Info",
				Rows: [][]string{
					{"Name", "Slug", "Route", "Method", "ExampleURL"},
				},
			},
			&admin.Section{
				Title: "Params",
				Rows: [][]string{
					{"Selector", "BaseURL", "PatternURL"},
				},
			},
			&admin.Section{
				Title: "Headers",
				Rows: [][]string{
					{"Headers"},
				},
			},
			&admin.Section{
				Title: "Blocks",
				Rows: [][]string{
					{"Blocks"},
				},
			},
			&admin.Section{
				Title: "Bots",
				Rows: [][]string{
					{"Extract"},
				},
			},
		)

		endpoint.EditAttrs(
			&admin.Section{
				Title: "Status",
				Rows: [][]string{
					{"Disabled", "Debug"},
				},
			},
			&admin.Section{
				Title: "Info",
				Rows: [][]string{
					{"Name", "Route", "Method", "ExampleURL"},
				},
			},
			&admin.Section{
				Title: "Params",
				Rows: [][]string{
					{"Selector", "BaseURL", "PatternURL"},
				},
			},
			&admin.Section{
				Title: "Headers",
				Rows: [][]string{
					{"Headers"},
				},
			},
			&admin.Section{
				Title: "Blocks",
				Rows: [][]string{
					{"Blocks"},
				},
			},
			&admin.Section{
				Title: "Bots",
				Rows: [][]string{
					{"Extract"},
				},
			},
		)
	*/
	endpointPropertiesRes := endpoint.Meta(&admin.Meta{Name: "EndpointProperties"}).Resource
	endpointPropertiesRes.NewAttrs(&admin.Section{
		Rows: [][]string{{"Name", "Value"}},
	})
	endpointPropertiesRes.EditAttrs(&admin.Section{
		Rows: [][]string{{"Name", "Value"}},
	})

	openapi := AdminUI.AddResource(&scraper.OpenAPIConfig{}, &admin.Config{Menu: []string{"API Scrapers"}})

	// Search resources
	// AdminUI.AddSearchResource(topic)
	AdminUI.AddSearchResource(endpoint)
	AdminUI.AddSearchResource(group)
	AdminUI.AddSearchResource(provider)
	AdminUI.AddSearchResource(openapi)

}

/*
	Refs:
	- https://github.com/dwarvesf/delivr-admin/blob/develop/config/admin/admin.go
	- https://github.com/xinuxZ/wzz_qor/blob/master/app/controllers/application.go
	- https://github.com/chenxin0723/ilove/blob/master/config/routes/routes.go
	- https://github.com/reechou/erp/blob/master/app/controllers/home.go
	- https://github.com/reechou/erp/blob/master/app/controllers/category.go
	- https://github.com/reechou/erp/blob/master/app/models/order.go
	- https://github.com/sunwukonga/paypal-qor-admin/blob/master/config/admin/admin.go
	- https://github.com/sunwukonga/paypal-qor-admin/blob/master/config/admin/admin.go
	- https://github.com/angeldm/optiqor/blob/master/app/controllers/application.go
	- https://github.com/angeldm/optiqor/blob/master/config/admin/admin.go
	- https://github.com/xinuxZ/wzz_qor/blob/master/config/admin/admin.go
	- https://github.com/sunwukonga/qor-scbn/blob/devmaster/config/admin/admin.go
	- https://github.com/damonchen/beezhu/blob/master/config/admin/admin.go
	- https://github.com/sunfmin/beego_with_qor/blob/master/main.go (beego+qor)

	- https://github.com/yalay/picCms/blob/dl/models/download.go
	- https://github.com/yalay/picCms/blob/dl/models/lang.go
	- https://github.com/8legd/hugocms/blob/master/qor/models/release.go
	- https://github.com/ROOT005/managesys/blob/master/models/client.go#L41
	- https://github.com/enlivengo/admincore/tree/master/views

	- https://github.com/ROOT005/com_web/blob/master/models/project.go
	- https://github.com/ROOT005/com_web/blob/master/models/website.go

	- https://github.com/drrzmr/pontalverde/blob/master/model/user/user.go
	- https://github.com/liujianping/scaffold/blob/master/templates/portal/app/models/base.db.go

*/
