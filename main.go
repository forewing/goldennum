package main

import (
	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/models"
	"github.com/forewing/goldennum/views"
	"github.com/gin-gonic/gin"
)

var (
	adminAcconts gin.Accounts = gin.Accounts{}
)

func main() {
	conf := config.Load()

	models.Load()
	defer models.Close()

	adminAcconts[conf.Admin] = conf.Password

	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	rAuth := r.Group("")
	// rAuth := r.Group("", gin.BasicAuth(adminAcconts))

	r.GET("/rooms", views.RoomList)
	r.GET("/room/:roomid", views.RoomInfo)

	rAuth.POST("/room", views.RoomCreate)
	rAuth.DELETE("/room/:roomid", views.RoomStop)
	rAuth.PUT("/room/:roomid", views.RoomStart)
	rAuth.PATCH("/room/:roomid", views.RoomUpdate)

	r.POST("/users/:roomid", views.UserCreate)

	r.GET("/user/:userid", views.UserInfo)
	r.POST("/user/:userid", views.UserSubmit)

	rAdmin := rAuth.Group("/admin")
	{
		rAdmin.GET("", views.AdminIndex)
	}

	if len(conf.Bind) == 0 {
		r.Run(":8080")
	} else {
		r.Run(conf.Bind)
	}
}
