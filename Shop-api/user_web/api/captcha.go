package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

// 默认的存储
var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	//生成driverDigit
	digit := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)

	//通过driver和store生成captcha
	captcha := base64Captcha.NewCaptcha(digit, store)
	captchaId, b64s, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码生成错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "生成验证码成功",
		"captchaId": captchaId,
		"b64s":      b64s,
	})
}
