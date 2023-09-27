package api

import (
	"Shop-api/user_web/forms"
	"Shop-api/user_web/global"
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(length int) string {
	number := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(uint64(time.Now().Unix()))
	var sb strings.Builder
	for i := 0; i < length; i++ {
		_, err := fmt.Fprintf(&sb, "%d", rand.Intn(len(number)))
		if err != nil {
			zap.S().Errorw("[SendSms.GenerateSmsCode] Error : ", err.Error())
			return ""
		}
	}
	return sb.String()
}

func SendSms(c *gin.Context) {
	//c.JSON()
	sendSmsForm := forms.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		zap.S().Errorw("[SendSms] Error : ", err.Error())
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-wuhan",
		"LTAI5tAnKRH5tFTvgtSqJedQ", "R1q4cqtVKRJN48Ud98RYWf7JKpIYye")
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile            //手机号
	request.QueryParams["SignName"] = "慕学在线"                            //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_181850725"               //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	if err != nil {
		fmt.Print(err.Error())
	}

	global.Rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, global.ExpiredTime)
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
