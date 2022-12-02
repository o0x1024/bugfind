package bugScan

import (
	"bufio"
	"bugfind/global"
	"bugfind/internal/utils"
	types2 "bugfind/model/types"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestRunBugScan(t *testing.T) {

	log.Println("[*] start bug scan...")
	cmd := exec.Command("nuclei", "-u", "https://rw.online.anaheim.cust66.lv.webproxy.ida.webank.com")

	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe: ", err)
		return
	}
	cmd.Stderr = os.Stderr
	//cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return
	}
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		data, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		data = strings.Replace(data, "\u001B[36m", "", -1)
		data = strings.Replace(data, "\u001B[1;92m", "", -1)
		data = strings.Replace(data, "\u001B[92m", "", -1)
		data = strings.Replace(data, "\u001B[94m", "", -1)
		data = strings.Replace(data, "\u001B[34m", "", -1)
		data = strings.Replace(data, "\u001B[33m", "", -1)
		data = strings.Replace(data, "\u001B[0m", "", -1)
		data = strings.Replace(data, "\u001B[38;5;208m", "", -1)

		for _, v := range keyword {
			reg := regexp.MustCompile(v)
			if len(reg.FindAllStringSubmatch(data, -1)) > 0 {
				//if strings.Contains(data, v) {
				fmt.Println(data)
				var alarm = types2.Alarm{}
				alarm.Type = 0
				alarm.Level = v
				alarm.Detail = data
				alarm.Timetamp = utils.GetCurTime()
				//扫描出来的漏洞上报后台
				global.Alarm <- alarm
				//扫描出来的漏洞微信公众号通知
			}
		}
	}
	err = cmd.Wait()

	log.Println("[*] bug scan done.")
}

func TestRunBugScans(t *testing.T) {

	var keywords = []string{`[.*(low.*]`, `\[.*medium.*\]`, `\[.*high.*\]`, `\[.*critical.*\]`}

	data := "[2022-03-15 04:19:39] [CVE-2016-6210] [network] [medium] fkysr2.iqiyi.com:22 [SSH-2.0-OpenSSH_5.3] [Hostname=fkysr2.iqiyi.com]"

	for _, v := range keywords {
		reg := regexp.MustCompile(v)
		match := reg.FindAllString(data, -1)
		if len(match) > 0 {
			//if strings.Contains(data, v) {
			fmt.Println(match)
		}
	}
}
