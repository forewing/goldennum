package main

import (
	"log"
	"os"

	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/models"
	"github.com/forewing/goldennum/views"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

const (
	staticsPath   = "./statics"
	templatesPath = "./templates"
)

var (
	staticsBox = packr.New("statics", staticsPath)

	adminAccounts gin.Accounts = gin.Accounts{}
)

func main() {
	conf := config.Load()

	models.Load()
	defer models.Close()

	adminAccounts[conf.Admin] = conf.Password

	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// pages
	if conf.Debug && canLiveReload(templatesPath) {
		log.Println("[main] templates use live reload")
		r.LoadHTMLGlob(templatesPath + "/*")
	} else {
		t, err := views.LoadTemplate()
		if err != nil {
			panic(err)
		}
		r.SetHTMLTemplate(t)
	}
	{
		r.GET("/", views.PageIndex)
	}

	// static files
	if conf.Debug && canLiveReload(staticsPath) {
		log.Println("[main] statics use live reload")
		r.Static("/statics", staticsPath)
	} else {
		r.StaticFS("/statics", staticsBox)
	}

	// public API
	{
		r.GET("/rooms", views.RoomList)
		r.GET("/room/:roomid", views.RoomInfo)
		r.GET("/sync/:roomid", views.RoomSync)

		r.POST("/users/:roomid", views.UserCreate)
		r.GET("/user/:userid", views.UserInfo)
		r.POST("/user/:userid", views.UserSubmit)
	}

	// admin API
	// rAdmin := r.Group("") // for test only
	rAdmin := r.Group("", gin.BasicAuth(adminAccounts))
	{
		rAdmin.GET("/admin", views.AdminIndex)
		rAdmin.POST("/room", views.RoomCreate)
		rAdmin.DELETE("/room/:roomid", views.RoomStop)
		rAdmin.PUT("/room/:roomid", views.RoomStart)
		rAdmin.PATCH("/room/:roomid", views.RoomUpdate)
	}

	if len(conf.Bind) == 0 {
		r.Run(":8080")
	} else {
		r.Run(conf.Bind)
	}
}

func canLiveReload(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
