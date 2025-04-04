//package utils
//
//import (
//	"encoding/json"
//	"fmt"
//	"net/http"
//)
//
//func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(statusCode)
//	err := json.NewEncoder(w).Encode(data)
//	if err != nil {
//		fmt.Fprintf(w, "%s", err.Error())
//	}
//}
//func ERROR(w http.ResponseWriter, statusCode int, err error) {
//	if err != nil {
//		JSON(w, statusCode, map[string]interface{}{
//			"message": err.Error(),
//		})
//		return
//	}
//	JSON(w, http.StatusBadRequest, nil)
//}

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Response(context *gin.Context, statusCode int, data interface{}) {
	context.Set("log-data", data)
	context.JSON(statusCode, data)
}

func SuccessResponse(context *gin.Context, data interface{}) {
	Response(context, http.StatusOK, data)
}

func ErrorResponse(context *gin.Context, statusCode int, message string) {
	Response(context, statusCode, Error{Code: statusCode, Message: message})
}
