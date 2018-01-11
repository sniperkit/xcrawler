package ajii

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//
// Frontend handlers
//
func Main(conf *EtcdConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{"foo": conf.foo},
		)
	}
}

func Frontend(port string, conf *EtcdConfig) {
	var router *gin.Engine
	gin.SetMode(gin.ReleaseMode)
	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.

	// router.LoadHTMLGlob("templates/*")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/", Main(conf))

	// Start serving the application
	router.Run(port)
}
