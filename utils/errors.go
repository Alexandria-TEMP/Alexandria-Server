package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Using http error handling as in:
// https://github.com/swaggo/swag/blob/master/example/celler/httputil/error.go

func HTTPError(ctx *gin.Context, status int, err error) {
	er := struct {
		Code    int
		Message string
	}{status, err.Error()}

	fmt.Println(err)
	ctx.JSON(status, er)
}
