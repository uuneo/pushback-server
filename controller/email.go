package controller

import (
	"NewBearService/config"
	"NewBearService/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type EmailVerificationCode struct {
	Code   string `json:"code"`
	Count  int    `json:"count"`
	Expire int64  `json:"expire"`
}

func (e *EmailVerificationCode) IsExpired() bool {
	return e.Expire < time.Now().Unix() || e.Count > 4
}

var EMAILVERIFICATIONCODEMAP = make(map[string]EmailVerificationCode)

func SendEmail(emailConfig *config.Email, toEmail string) (bool, string) {

	generateVerificationCode := func() string {
		// 设置随机数种子
		rand.NewSource(time.Now().UnixNano())
		// 生成 6 位随机数
		randomNumber := rand.Intn(900000) + 100000
		// 替换验证码中的 4 为 5
		randomNumber1 := strings.Replace(fmt.Sprintf("%d", randomNumber), "4", "7", -1)
		return randomNumber1
	}

	randomCode := generateVerificationCode()

	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", emailConfig.Subject)
	m.SetBody("text/html", randomCode+emailConfig.BodySuffix)

	d := gomail.NewDialer(emailConfig.SmtpHost, emailConfig.SmtpPort, emailConfig.From, emailConfig.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return false, "发送失败"
	}

	EMAILVERIFICATIONCODEMAP[toEmail] = EmailVerificationCode{
		Code:   randomCode,
		Count:  0,
		Expire: time.Now().Add(time.Minute * 3).Unix(),
	}

	go func() {
		for key, v := range EMAILVERIFICATIONCODEMAP {
			if v.IsExpired() {
				delete(EMAILVERIFICATIONCODEMAP, key)
			}
		}
	}()

	return true, randomCode
}

func SendCode(c *gin.Context) {
	email := c.Query("email")
	if datap, ok := EMAILVERIFICATIONCODEMAP[email]; ok {
		if !datap.IsExpired() {
			c.JSON(200,
				failed(400,
					fmt.Sprintf(
						"expire time: %v",
						datap.Expire-time.Now().Unix(),
					),
				),
			)
			return
		}

	}
	if email == "" {
		c.JSON(200, failed(400, "email is required"))
		return
	}
	emailConfig := config.LocalConfig.Email
	suc, code := SendEmail(&emailConfig, email)
	fmt.Println("send email code: ", code)

	if suc {
		c.JSON(200, success())
	} else {
		c.JSON(200, failed(500, "send email failed"))
	}
}

func VerifyCode(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")
	key := c.Query("key")
	deviceToken := c.Query("deviceToken")
	if email == "" || code == "" || deviceToken == "" {
		c.JSON(200, failed(400, "email, code, deviceToken is required"))
		return
	}

	if v, ok := EMAILVERIFICATIONCODEMAP[email]; ok {
		if v.IsExpired() {
			c.JSON(200, failed(400, "code is expired"))
			delete(EMAILVERIFICATIONCODEMAP, email)
			return
		}
		if v.Code == code {
			newKey, err := database.DB.SaveDeviceTokenByEmail(email, key, deviceToken)

			if err != nil {
				c.JSON(http.StatusOK, failed(500, "device registration failed: %v", err))
				return
			}

			c.JSON(http.StatusOK, data(map[string]string{
				"key":          newKey,
				"device_key":   newKey,
				"device_token": deviceToken,
			}))
			delete(EMAILVERIFICATIONCODEMAP, email)
			return
		} else {
			v.Count++
			EMAILVERIFICATIONCODEMAP[email] = v
		}
	}

	c.JSON(200, failed(400, "code is wrong"))
}
