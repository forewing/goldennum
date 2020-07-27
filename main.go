package main

import (
	"html/template"

	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/models"
	"github.com/forewing/goldennum/views"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

var (
	staticsBox   = packr.New("statics", "./statics")
	templatesBox = packr.New("templates", "./templates")

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
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
	{
		r.GET("/", views.PageIndex)
	}

	// static files
	{
		r.StaticFS("/statics", staticsBox)
	}

	// public API
	{
		r.GET("/rooms", views.RoomList)
		r.GET("/room/:roomid", views.RoomInfo)

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

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range views.Templates {
		str, err := templatesBox.FindString(name)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(str))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
