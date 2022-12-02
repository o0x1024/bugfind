package taskCore

import (
	"BugFind/global"
	"BugFind/internal/assetCollect/domain"
	"BugFind/model/response"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func Run() {
	timer := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timer.C:
			task, err := getTask()
			if err == nil && task != nil {
				execTask(task)
			}
		}
	}
}

// taskid
// 1000 信息搜集-域名-全部方式
// 1001 信息搜集-域名-接口扫描
// 1002 信息搜集-域名-暴力破解
// 1003 信息搜集-域名-暴力破解
func execTask(task *response.Task) {
	log.Println("found task task Id:", task.TaskId)
	switch task.TaskId {
	case 1000:
	case 1001:
		domain.ColletDomainByInterface(task)
	case 1002:
		domain.ColletDomainByKSubdomain(task)
	case 1003:

	case 1004:

	}
}

func getTask() (task *response.Task, err error) {
	GetTaskUrl := global.WkgURL + "/v3/task/getTask?token=" + global.V3Token
	client := http.Client{}

	req, _ := http.NewRequest("GET", GetTaskUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	wr := response.WkgRes{}
	err = json.Unmarshal(buf, &wr)
	if err != nil {
		log.Println(err)
		return
	}

	if wr.Code == 200 {
		return &wr.Task, nil
	}

	return
}
