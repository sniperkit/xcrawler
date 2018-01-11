package ajii

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	// "time"
)

//
// Backend shit
//
type CallCtx struct {
	Context *gin.Context
}

func createTask(ctx *CallCtx, callback func(*CallCtx) interface{}) (func(), int) {
	fmt.Println("Spawning a thing")
	return func() {
		fmt.Println("Running a thing")
		callback(ctx)
		fmt.Println("Done")
	}, 0
}

func ConfigMgr(conf *EtcdConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		conf.foo = "biz"
		c.JSON(
			http.StatusOK,
			gin.H{"hello": "world"},
		)
	}
}

func PostTask(conf *EtcdConfig, callback func(*CallCtx) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &CallCtx{Context: c}
		// task, task_uid := createTask(ctx, callback)
		// go task()
		rp := callback(ctx)
		c.JSON(
			http.StatusOK,
			rp,
			// gin.H{"accepted:": task_uid},
		)
	}

}

func GetTaskLst(conf *EtcdConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []string
		c.JSON(
			http.StatusOK,
			gin.H{"tasks": tasks},
		)
	}
}

func Backend(port string, conf *EtcdConfig, callback func(*CallCtx) interface{}) {
	var router *gin.Engine
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.GET("/", ConfigMgr(conf))
	router.POST("/v1/task", PostTask(conf, callback))
	router.GET("/v1/task", GetTaskLst(conf))
	router.Run(port)
}
