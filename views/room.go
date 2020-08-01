package views

import (
	"fmt"
	"log"
	"net/http"

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
		log.Printf("Info: [views] getRoomByIDOrErr, caller: %v, ID: %v, %v\n", caller, roomid, result.Error)
		c.JSON(http.StatusNotFound, "")
		return nil, result.Error
	} else if result.Error != nil {
		log.Printf("Error: [views] getRoomByIDOrErr, ID: %v, %v\n", roomid, result.Error)
		c.JSON(http.StatusInternalServerError, "")
		return nil, result.Error
	}

	return &room, nil
}

// RoomCreate create room
func RoomCreate(c *gin.Context) {
	var data roomCreateModel
	if err := c.BindJSON(&data); err != nil {
		log.Printf("Info: [views] RoomCreate, %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	room := models.Room{
		Interval:   data.Interval,
		RoundNow:   0,
		RoundTotal: data.RoundTotal,
	}
	if err := models.Db.Create(&room).Error; err != nil {
		log.Printf("Error: [views] RoomCreate, %v\n", err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	log.Printf("Info: [views] RoomCreate, create room, %+v\n", room)

	room.Start()

	// user, _ := models.UserNew(room.ID, fmt.Sprintf("u%v", rand.Uint32()), "12345678")
	// models.Db.Save(user)

	c.JSON(http.StatusOK, room)
}

// RoomList list all rooms
func RoomList(c *gin.Context) {
	rooms := []models.Room{}
	if err := models.Db.Find(&rooms).Error; err != nil {
		log.Printf("Error: [views] RoomList, %v\n", err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	log.Printf("Info: [views] RoomList, len: %v\n", len(rooms))
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
		return
	}
	duration, err := models.RoomUntilNextTick(uint(roomid))
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("%.0f", duration.Seconds()))
}
