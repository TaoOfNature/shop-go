package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Envelope{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func Error(c *gin.Context, status int, code int, msg string) {
	c.JSON(status, Envelope{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
