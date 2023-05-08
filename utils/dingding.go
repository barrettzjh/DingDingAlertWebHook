package utils

import DingDing "github.com/JetBlink/dingtalk-notify-go-sdk"

func SendDingTalk(webHook, secret string, content interface{}) error {
	return DingDing.NewRobot(webHook, secret).SendMessage(content)
}
