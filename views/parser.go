package views

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func parseInt64FromParamOrErr(c *gin.Context, key string, caller string) (int64, error) {
	num, err := strconv.ParseInt(c.Param(key), 10, 0)
	if err != nil {
		zap.S().Warnf("%v->ParseInt64FromParamOrErr, %v", caller, err)
		c.JSON(http.StatusBadRequest, err.Error())
		return -1, err
	}
	return num, nil
}
