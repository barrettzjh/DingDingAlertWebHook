package utils

import dingtalk_robot "github.com/JetBlink/dingtalk-notify-go-sdk"

func SendDingTalkResolved(content interface{}) error {
	//webHook := os.Getenv("WEBHOOK")
	webHook := "a354dd86ae189294a7ea49746beed6ab74a0247036a7a525c60c6e8ce5846168"
	robot := dingtalk_robot.NewRobot(webHook, "ceshi")
	//return robot.SendTextMessage(content, []string{}, false)

	return robot.SendMessage(content)
}

func SendDingTalkResolvedWithWebHook(webHook, secret string, content interface{}) error {
	robot := dingtalk_robot.NewRobot(webHook, secret)
	//return robot.SendTextMessage(content, []string{}, false)

	return robot.SendMessage(content)
}
