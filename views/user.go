package views

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/forewing/goldennum/models"

	"github.com/forewing/goldennum/utils"
	"github.com/gin-gonic/gin"
)

func getUserByIDOrErr(userid int64, c *gin.Context, caller string) (*models.User, error) {
	var user models.User
	if result := models.Db.First(&user, userid); result.RecordNotFound() {
		zap.S().Warnf("getUserByIDOrErr, caller: %v, ID: %v, %v", caller, userid, result.Error)
		c.JSON(http.StatusNotFound, fmt.Sprintf("User %v", userid))
		return nil, result.Error
	} else if result.Error != nil {
		zap.S().Errorf("getUserByIDOrErr, ID: %v, %v", userid, result.Error)
		c.JSON(http.StatusInternalServerError, "")
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
		zap.S().Warnf("UserCreate, %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if !models.UserNameValidate(data.Username) {
		zap.S().Warnf("UserCreate, invalid username: %v", data.Username)
		c.JSON(http.StatusBadRequest, "invalid username")
	}
	if !models.UserPassValidate(data.Password) {
		zap.S().Warnf("UserCreate, invalid password, len: %v", len(data.Password))
		c.JSON(http.StatusBadRequest, "invalid password")
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
		c.JSON(http.StatusInternalServerError, "fail to create user")
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
	// if user.Auth()
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
		zap.S().Warnf("UserSubmit, %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := user.Auth(data.Password); err != nil {
		zap.S().Warnf("UserSubmit, auth: %v", err)
		c.JSON(http.StatusUnauthorized, "")
		return
	}
	if !models.UserSubmitValidate(data.Submit1) || !models.UserSubmitValidate(data.Submit2) {
		zap.S().Warnf("UserSubmit, submit: %v, %v", data.Submit1, data.Submit2)
		c.JSON(http.StatusBadRequest, "")
		return
	}
	if err := user.Submit(data.Submit1, data.Submit2); err != nil {
		zap.S().Warnf("UserSubmit, submit: %v", err)
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	zap.S().Infof("UserSubmit, %v, %v", data.Submit1, data.Submit2)
	c.JSON(http.StatusOK, "")
}

type userAuthModel struct {
	Password string `json:"Password" binding:"required"`
}

// UserAuth check user credential
func UserAuth(c *gin.Context) {
	userid, err := utils.ParseInt64FromParamOrErr(c, "userid", "UserAuth")
	if err != nil {
		return
	}
	var data userAuthModel
	if err := c.BindJSON(&data); err != nil {
		zap.S().Warnf("UserAuth, %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := getUserByIDOrErr(userid, c, "UserAuth")
	if err != nil {
		return
	}
	if err := user.Auth(data.Password); err != nil {
		zap.S().Warnf("UserAuth, auth: %v", err)
		c.JSON(http.StatusUnauthorized, "")
		return
	}
	c.JSON(http.StatusOK, user)
}
