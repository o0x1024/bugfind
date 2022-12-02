package notice

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestProgramErrorAlarm(t *testing.T) {
	webhook := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=71a8d6f3-f15a-4a43-a323-5e411ae3e9b9"

	level := "low"
	Detail := "Get https://fofa.info/api/v1/search/all?email=zjgelen@gmail.com&key=1bd13cc61d22823099fea2a8e26f7478&qbase64=Y2VydD0iYWkucGguY29tLmNuIg==&page=1&size=10000&fields=host"
	timetamp := "2022-2-24 12:44:11"
	content := fmt.Sprintf(`发现新漏洞<font color=\"warning\">1个</font>，请注意。\n
        >级别:<font color=\"comment\">%s</font>
        >细节:<font color=\"comment\">%s</font>
        >时间:<font color=\"comment\">%s</font>`,level,Detail,timetamp)


	msg := strings.Replace(`{"msgtype": "markdown","markdown": {"content": "_gelen_"}}`,"_gelen_",content,-1)


	req ,_:= http.NewRequest("POST",webhook,strings.NewReader(msg))
	req.Header.Add("Content-Type","application/json;charset=utf-8")

	client := http.Client{Timeout: 3*time.Second}
	_ ,err := client.Do(req)
	if err != nil{
		log.Println("webchat notice error[",err,"]")
		return
	}
}
