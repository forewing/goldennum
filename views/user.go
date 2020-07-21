package views

import (
	"log"
	"net/http"

	"github.com/forewing/goldennum/models"

	"github.com/forewing/goldennum/utils"
	"github.com/gin-gonic/gin"
)

func getUserByIDOrErr(userid int64, c *gin.Context, caller string) (*models.User, error) {
	var user models.User
	if result := models.Db.First(&user, userid); result.RecordNotFound() {
		log.Printf("Info: [views] getUserByIDOrErr, caller: %v, ID: %v, %v\n", caller, userid, result.Error)
		c.String(http.StatusNotFound, "")
		return nil, result.Error
	} else if result.Error != nil {
		log.Printf("Error: [views] getUserByIDOrErr, ID: %v, %v\n", userid, result.Error)
		c.String(http.StatusInternalServerError, "")
		return nil, result.Error
	}

	return &user, nil
}

type userCreateModel struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

// UserCreate create user
func UserCreate(c *gin.Context) {
	var data userCreateModel
	if err := c.BindJSON(&data); err != nil {
		log.Printf("Info: [views] UserCreate, %v\n", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if !models.UserNameValidate(data.Username) {
		log.Printf("Info: [views] UserCreate, invalid username: %v\n", data.Username)
		c.String(http.StatusBadRequest, "invalid username")
	}
	if !models.UserPassValidate(data.Password) {
		log.Printf("Info: [views] UserCreate, invalid password, len: %v\n", len(data.Password))
		c.String(http.StatusBadRequest, "invalid password")
	}

	roomid, err := utils.ParseInt64FromParamOrErr(c, "roomid", "UserCreate")
	if err != nil {
		return
	}
	room, err := getRoomByIDOrErr(roomid, c, "UserCreate")
	if err != nil {
		return
	}

	user, err := models.UserNew(room.ID, data.Username, data.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "fail to create user")
		return
	}
	models.Db.Save(&user)
	c.JSON(http.StatusOK, user)
}

// UserInfo show user info
func UserInfo(c *gin.Context) {
	userid, err := utils.ParseInt64FromParamOrErr(c, "userid", "UserInfo")
	if err != nil {
		return
	}
	user, err := getUserByIDOrErr(userid, c, "UserInfo")
	if err != nil {
		return
	}
	user.FilterInfo(false)
	user.UserHistorys = user.GetHistory()
	c.JSON(http.StatusOK, user)
}

type userSubmitModel struct {
	Password string  `json:"Password" binding:"required"`
	Submit1  float64 `json:"Submit1" binding:"required"`
	Submit2  float64 `json:"Submit2" binding:"required"`
}

// UserSubmit submit user action
func UserSubmit(c *gin.Context) {
	userid, err := utils.ParseInt64FromParamOrErr(c, "userid", "UserSubmit")
	if err != nil {
		return
	}
	user, err := getUserByIDOrErr(userid, c, "UserSubmit")
	if err != nil {
		return
	}
	var data userSubmitModel
	if err := c.BindJSON(&data); err != nil {
		log.Printf("Info: [views] UserSubmit, %v\n", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := user.Auth(data.Password); err != nil {
		log.Printf("Info: [views] UserSubmit, auth: %v\n", err)
		c.String(http.StatusUnauthorized, "")
		return
	}
	if !models.UserSubmitValidate(data.Submit1) || !models.UserSubmitValidate(data.Submit2) {
		log.Printf("Info: [views] UserSubmit, submit: %v, %v\n", data.Submit1, data.Submit2)
		c.String(http.StatusBadRequest, "")
		return
	}
	user.Submit1 = data.Submit1
	user.Submit2 = data.Submit2
	if result := models.Db.Save(&user); result.Error != nil {
		log.Printf("Error: [views] UserSubmit, save: %v\n", result.Error)
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.String(http.StatusOK, "")
}
