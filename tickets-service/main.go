package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"train-tickets-service/tools"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

func main() {
	// 实现扫码登录
	client := resty.New()
	getQrResp, getQrErr := tools.Post(context.TODO(), client, "https://kyfw.12306.cn/passport/web/create-qr64", nil, map[string]string{
		"appid": "otn",
	}, nil, true)
	if getQrErr != nil {
		log.Fatalln(getQrErr)
		return
	}

	base64Img := gjson.Get(getQrResp.String(), "image").String()
	BrowserQr(base64Img)
	uuid := gjson.Get(getQrResp.String(), "uuid").String()

	// 轮训监听二维码
	for {
		time.Sleep(time.Second * 1)
		resp, err := tools.Post(context.TODO(), client, "https://kyfw.12306.cn/passport/web/checkqr", nil, map[string]string{
			"RAIL_DEVICEID":   "",
			"RAIL_EXPIRATION": "",
			"uuid":            uuid,
			"appid":           "otn",
		}, nil, false)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		if gjson.Get(resp.String(), "result_code").String() == "2" {
			fmt.Println(gjson.Get(resp.String(), "uamtk").String())
			break
		}
		if gjson.Get(resp.String(), "result_code").String() == "3" {
			log.Fatalln(gjson.Get(resp.String(), "result_message").String())
			break
		}
	}
}

func BrowserQr(base64Image string) (err error) {
	var imgData []byte
	if imgData, err = base64.StdEncoding.DecodeString(base64Image); err != nil {
		return
	}
	imgPath := "./qr.png"
	err = ioutil.WriteFile(imgPath, imgData, 0644)
	if err != nil {
		fmt.Println("无法保存图片:", err)
		return
	}
	return
}

func GetQr() {

}
