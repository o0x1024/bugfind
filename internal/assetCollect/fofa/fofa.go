package fofa

import (
	"bugfind/global"
	utils2 "bugfind/internal/utils"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/kataras/golog"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	tnum     = 10
	wg       = &sync.WaitGroup{}
	ruleChan = make(chan string)
)

func CollectByFofa(ruleList []string) {
	log.Println("[*] fofa start.")
	wg.Add(tnum)
	for i := 0; i < tnum; i++ {
		go work()
	}
	go HandlerParam(ruleList)

	wg.Wait()

	log.Println("[+] fofa done.")
	global.FofaDone <- true
}

func HandlerParam(Params []string) {
	for _, v := range Params {
		ruleChan <- v
	}
	close(ruleChan)
}

func work() {
	var baseUrl string
	var fofaKey string
	var fofaEmail string
	var rule string
	var ok bool
	var err error
	var resp *http.Response

	fofaEmail = "zjgelen@gmail.com"
	fofaKey = "1bd13cc61d22823099fea2a8e26f7478"

	rule, ok = <-ruleChan

	for ok {
		//随机生成1-5秒休眠时间
		time.Sleep(time.Duration(utils2.RandomGenSleepTime()) * time.Second)
		_rule_b64str := base64.StdEncoding.EncodeToString([]byte(rule))

		baseUrl = fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&key=%s&qbase64=%s&page=1&size=10000&fields=host", fofaEmail, fofaKey, _rule_b64str)

		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := http.Client{Timeout: 5 * time.Second, Transport: tr}

		resp, err = client.Get(baseUrl)
		if err != nil {
			golog.Error("fofa.go line:30 ", err)
			global.Alarm <- utils2.GenReportInfo(err)
			rule, ok = <-ruleChan
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			golog.Error("fofa.go line:36 ", err)
			global.Alarm <- utils2.GenReportInfo(err)
			rule, ok = <-ruleChan
			continue
		}

		jsdata, _ := simplejson.NewJson(body)
		results, err := jsdata.Get("results").Array()
		if err != nil {
			s, _ := jsdata.Get("errmsg").String()
			log.Println("[err] from fofa response " + s)
			global.Alarm <- utils2.GenReportInfo(err)
			rule, ok = <-ruleChan
			continue
		}
		log.Println("[+] find ", len(results), "assets by fofa.")
		for _, v := range results {
			switch v.(type) {
			case string:
				global.Target <- v.(string)
			}
		}
		if resp != nil {
			resp.Body.Close()
		}
		rule, ok = <-ruleChan
	}

	wg.Done()
}
