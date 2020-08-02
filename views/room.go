package views

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/forewing/goldennum/models"
	"github.com/forewing/goldennum/utils"
	"github.com/gin-gonic/gin"
)

type roomCreateModel struct {
	Interval   int `json:"Interval" binding:"required"`
	RoundTotal int `json:"RoundTotal" binding:"required"`
}

func getRoomByIDOrErr(roomid int64, c *gin.Context, caller string) (*models.Room, error) {
	var room models.Room
	if result := models.Db.First(&room, roomid); result.RecordNotFound() {
		zap.S().Warnf("getRoomByIDOrErr, caller: %v, ID: %v, %v", caller, roomid, result.Error)
		c.JSON(http.StatusNotFound, "")
		return nil, result.Error
	} else if result.Error != nil {
		zap.S().Errorf("getRoomByIDOrErr, ID: %v, %v", roomid, result.Error)
		c.JSON(http.StatusInternalServerError, "")
		return nil, result.Error
	}

	return &room, nil
}

// RoomCreate create room
func RoomCreate(c *gin.Context) {
	var data roomCreateModel
	if err := c.BindJSON(&data); err != nil {
		zap.S().Warnf("RoomCreate, %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	room := models.Room{
		Interval:   data.Interval,
		RoundNow:   0,
		RoundTotal: data.RoundTotal,
	}
	if err := models.Db.Create(&room).Error; err != nil {
		zap.S().Errorf("RoomCreate, %v", err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	zap.S().Infof("RoomCreate, create room, %+v", room)

	room.Start()

	// user, _ := models.UserNew(room.ID, fmt.Sprintf("u%v", rand.Uint32()), "12345678")
	// models.Db.Save(user)

	c.JSON(http.StatusOK, room)
}

// RoomList list all rooms
func RoomList(c *gin.Context) {
	rooms := []models.Room{}
	if err := models.Db.Find(&rooms).Error; err != nil {
		zap.S().Errorf("RoomList, %v", err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	zap.S().Infof("RoomList, len: %v", len(rooms))
	c.JSON(http.StatusOK, rooms)
}

// RoomInfo show room info
func RoomInfo(c *gin.Context) {
	roomid, err := utils.ParseInt64FromParamOrErr(c, "roomid", "RoomInfo")
	if err != nil {
		return
	}
	room, err := getRoomByIDOrErr(roomid, c, "RoomInfo")
	if err != nil {
		return
	}
	room.Users = room.GetUsers()
	room.RoomHistorys = room.GetHistory()
	c.JSON(http.StatusOK, room)
}

// RoomStart start the room
func RoomStart(c *gin.Context) {
	roomid, err := utils.ParseInt64FromParamOrErr(c, "roomid", "RoomInfo")
	if err != nil {
		return
	}
	room, err := getRoomByIDOrErr(roomid, c, "RoomInfo")
	if err != nil {
		return
	}
	if room.Start() {
		c.JSON(http.StatusOK, "")
	} else {
		c.JSON(http.StatusBadRequest, "")
	}
}

// RoomStop stop the room
func RoomStop(c *gin.Context) {
	roomid, err := utils.ParseInt64FromParamOrErr(c, "roomid", "RoomInfo")
	if err != nil {
		return
	}
	room, err := getRoomByIDOrErr(roomid, c, "RoomInfo")
	if err != nil {
		return
	}
	if room.Stop() {
		c.JSON(http.StatusOK, "")
	} else {
		c.JSON(http.StatusBadRequest, "")
	}
}

// RoomUpdate update the room
func RoomUpdate(c *gin.Context) {

}

// RoomSync return time until next tick
func RoomSync(c *gin.Context) {
	roomid, err := utils.ParseInt64FromParamOrErr(c, "roomid", "RoomInfo")
	if err != nil {
		zap.S().Warnf("RoomSync, %v", err)
		return
	}
	duration, err := models.RoomUntilNextTick(uint(roomid))
	if err != nil {
		zap.S().Warnf("RoomSync, %v", err)
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	zap.S().Infof("RoomSync, %v", duration.String())
	c.JSON(http.StatusOK, fmt.Sprintf("%.0f", duration.Seconds()))
}
