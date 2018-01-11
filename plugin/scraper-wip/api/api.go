package api

/*
import (
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/roscopecoltran/scraper/db"
	"github.com/roscopecoltran/scraper/scraper"
)
*/

/*
	Refs:
	"github.com/wantedly/webmock-proxy"
	"github.com/wantedly/apig"
	"github.com/Gacnt/gormid"
	"github.com/huytd/api-togo"
*/

//var API *admin.Admin

//func init() {
//	API = admin.New(&qor.Config{DB: db.DB})

/*
	Product := API.AddResource(&models.Product{})

	ColorVariationMeta := Product.Meta(&admin.Meta{Name: "ColorVariations"})
	ColorVariation := ColorVariationMeta.Resource
	ColorVariation.IndexAttrs("ID", "Color", "Images", "SizeVariations")
	ColorVariation.ShowAttrs("Color", "Images", "SizeVariations")

	SizeVariationMeta := ColorVariation.Meta(&admin.Meta{Name: "SizeVariations"})
	SizeVariation := SizeVariationMeta.Resource
	SizeVariation.IndexAttrs("ID", "Size", "AvailableQuantity")
	SizeVariation.ShowAttrs("ID", "Size", "AvailableQuantity")

	API.AddResource(&models.Order{})

	User := API.AddResource(&models.User{})
	userOrders, _ := User.AddSubResource("Orders")
	userOrders.AddSubResource("OrderItems", &admin.Config{Name: "Items"})

	API.AddResource(&models.Category{})
*/
//}
