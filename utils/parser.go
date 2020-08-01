package utils

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseInt64FromParamOrErr(c *gin.Context, key string, caller string) (int64, error) {
	num, err := strconv.ParseInt(c.Param(key), 10, 0)
	if err != nil {
		log.Printf("Info: [utils] %v-ParseInt64FromParamOrErr, %v\n", caller, err)
		c.JSON(http.StatusBadRequest, err.Error())
		return -1, err
	}
	return num, nil
}
