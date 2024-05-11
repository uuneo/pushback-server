package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func QRCode(c *gin.Context) {
	url := "https://push.twown.com"
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		c.JSON(200, "生成二维码失败")
		return
	}
	c.Writer.Write(png)

}
