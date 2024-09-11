package main

import (
	"fmt"

	"github.com/PirateDreamer/going/conf"
	"github.com/eatmoreapple/openwechat"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var UIDMapRemark = make(map[string]string)

func main() {
	conf.InitConfig(nil)
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数c
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() {
			fromUser, _ := msg.Sender()
			if fromUser.RemarkName == viper.GetString("ai_reply_remark_name") {
				resp, err := resty.New().R().SetBody(map[string]any{
					"question": msg.Content,
				}).Post("http://localhost:8080/api/ollama/get_resp")
				if err != nil {
					return
				}
				if resp.StatusCode() != 200 {
					return
				}
				if gjson.Get(resp.String(), "code").Int() != 0 {
					return
				}
				msg.ReplyText(gjson.Get(resp.String(), "data.reply").String())
			}
		}
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
