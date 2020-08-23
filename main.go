package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/models"
	"github.com/forewing/goldennum/views"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	staticsPath   = "statics"
	templatesPath = "templates"
)

var (
	adminAccounts gin.Accounts = gin.Accounts{}
)

//go:generate go run generate/main.go
func main() {
	defer setLogger()()

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
		zap.S().Debugf("templates use live reload")
		r.LoadHTMLGlob(templatesPath + "/*")
	} else {
		t, err := loadTemplate()
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
		zap.S().Debugf("statics use live reload")
		r.Static("/statics", staticsPath)
	} else {
		r.StaticFS("/statics", AssetFile())
	}

	// public API
	{
		r.GET("/rooms", views.RoomList)
		r.GET("/room/:roomid", views.RoomInfo)
		r.GET("/sync/:roomid", views.RoomSync)

		r.POST("/users/:roomid", views.UserCreate)
		r.GET("/user/:userid", views.UserInfo)
		r.POST("/user/:userid", views.UserSubmit)
		r.PUT("/user/:userid", views.UserAuth)
	}

	// admin API
	// rAdmin := r.Group("") // for test only
	rAdmin := r.Group("", gin.BasicAuth(adminAccounts))
	r.GET("/admin", views.AdminIndex)
	{
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

func setLogger() func() error {
	conf := zap.NewDevelopmentConfig()
	logger, err := conf.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger.Sync
}

// LoadTemplate reutrn templates
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range views.Templates {
		data, err := Asset(templatesPath + "/" + name)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(data))
		if err != nil {
			return nil, err
		}
	}
	return t.New(views.TemplateBaseURL).Parse(
		fmt.Sprintf("{{ define \"%v\" }}%v{{ end }}", views.TemplateBaseURL, config.Load().BaseURL),
	)
}
