package main

/*

import (
    "fmt"
    "net/http"

	// Qor
    "github.com/qor/qor"
	"github.com/qor/l10n"
	"github.com/qor/media"
	"github.com/qor/publish2"
	"github.com/qor/sorting"
	"github.com/qor/validations"

	// Gorm / RDB
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	// Qor Admin
    "github.com/qor/admin"
)

// Create a GORM-backend model
type User struct {
  gorm.Model
  Name string
}

// Create another GORM-backend model
type Product struct {
  gorm.Model
  Name        string
  Description string
}

func main() {
  DB, _ := gorm.Open("sqlite3", "demo.db")
  DB.AutoMigrate(&User{}, &Product{})

  // Initalize
  Admin := admin.New(&qor.Config{DB: DB})

  // Create resources from GORM-backend model
  Admin.AddResource(&User{})
  Admin.AddResource(&Product{})

  // Register route
  mux := http.NewServeMux()
  // amount to /admin, so visit `/admin` to view the admin interface
  Admin.MountTo("/admin", mux)

  fmt.Println("Listening on: 9000")
  http.ListenAndServe(":9000", mux)
}

*/

/*
	// Beego
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)

	beego.Handler("/admin/*", mux)
	beego.Run()
*/

/*
	// Gin
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)

	r := gin.Default()
	r.Any("/admin/*w", gin.WrapH(mux))
	r.Run()
*/
