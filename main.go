package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/models"
	"github.com/forewing/goldennum/resources"
	"github.com/forewing/goldennum/views"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	staticsPath     = "statics"
	templatesPath   = "templates"
	templateBaseURL = "base_url"
)

var (
	adminAccounts gin.Accounts = gin.Accounts{}
)

func main() {
	defer setLogger()()

	conf := config.Load()

	models.Load()
	defer models.Close()

	adminAccounts[conf.Admin] = conf.Password

	if conf.Debug {
		resources.SetLiveReload()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// pages
	{
		t := mustLoadTemplate()
		r.SetHTMLTemplate(t)
		r.GET("/", views.PageIndex)
		r.GET("/admin", views.AdminIndex)
	}

	// static files
	{
		r.StaticFS("/statics", http.FS(resources.GetStatics()))
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
	{
		rAdmin.POST("/room", views.RoomCreate)
		rAdmin.DELETE("/room/:roomid", views.RoomStop)
		rAdmin.PUT("/room/:roomid", views.RoomStart)
		rAdmin.PATCH("/room/:roomid", views.RoomUpdate)
	}

	if len(conf.Bind) == 0 {
		zap.S().Info("Listening on http://127.0.0.1:8080")
		r.Run(":8080")
	} else {
		zap.S().Info("Listening on http://", conf.Bind)
		r.Run(conf.Bind)
	}
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

func mustLoadTemplate() *template.Template {
	templates := resources.GetTemplates()

	// load general templates
	t, err := template.ParseFS(templates, "*.html")
	if err != nil {
		panic(err)
	}

	// generate BaseURL template
	t, err = t.New(templateBaseURL).Parse(
		fmt.Sprintf("{{ define \"%v\" }}%v{{ end }}", templateBaseURL, config.Load().BaseURL),
	)
	if err != nil {
		panic(err)
	}
	return t
}
