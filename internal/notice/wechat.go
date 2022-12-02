package notice

import (
	"bugfind/config"
	"bugfind/global"
	"bugfind/internal/notice/request"
	types2 "bugfind/model/types"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// _type 0代表漏洞告警，1代表程序错误告警
func alarm(alarm types2.Alarm) {

	vulnWebhook := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + config.Cfg.Wechat
	expWebhook := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=040bd9c6-f2af-4a2b-8f56-65a20b6f505a"
	level := alarm.Level
	detail := alarm.Detail
	detail = strings.Replace(detail, "\"", "", -1)
	stack := alarm.Stack
	timetamp := alarm.Timetamp

	var levelColor string
	switch alarm.Level {
	case "[low]":
		levelColor = "orange"
	case "[medium]":
		levelColor = "warning"
	case "[high]":
		levelColor = "red"
	case "[critical]":
		levelColor = "red"
	}
	if alarm.Type == 1 {
		content := fmt.Sprintf(`<font color=\"warning\">程序运行异常，请注意</font>。\n
        >细节:<font color=\"comment\">%s</font>
		>堆栈:<font color=\"comment\">%s</font>
        >时间:<font color=\"comment\">%s</font>`, detail, stack, timetamp)

		msg := strings.Replace(`{"msgtype": "markdown","markdown": {"content": "_gelen_"}}`, "_gelen_", content, -1)

		req, _ := http.NewRequest("POST", expWebhook, strings.NewReader(msg))
		req.Header.Add("Content-Type", "application/json;charset=utf-8")

		client := http.Client{Timeout: 3 * time.Second}
		_, err := client.Do(req)
		if err != nil {
			log.Println("webchat notice error[", err, "]")
			return
		}
	} else if alarm.Type == 0 {

		content := fmt.Sprintf(`发现新漏洞<font color=\"warning\">1个</font>，请注意。\n
        >级别: <font color=\"%s\"> %s</font>
        >细节: <font color=\"comment\">%s</font>
        >时间: <font color=\"comment\"> %s</font>`, levelColor, level, detail, timetamp)

		msg := strings.Replace(`{"msgtype": "markdown","markdown": {"content": "_gelen_"}}`, "_gelen_", content, -1)

		req, _ := http.NewRequest("POST", vulnWebhook, strings.NewReader(msg))
		req.Header.Add("Content-Type", "application/json;charset=utf-8")

		client := http.Client{Timeout: 3 * time.Second}
		_, err := client.Do(req)
		if err != nil {
			log.Println("webchat notice error[", err, "]")
			return
		}
	}
}

func AlarmNoticeAndReport() {
	//程序告警通知
	go func() {
		for {
			select {
			case arm := <-global.Alarm:
				time.Sleep(5 * time.Second)
				//微信告警
				go alarm(arm)
				if arm.Type != 1 {
					//上传后端
					go request.BugReport(arm.Detail)
				}
			}
		}
	}()
}
