package handle

import (
	"bugfind/global"
	"bugfind/internal/libs/Glog"
	"bugfind/model/request"
	"bugfind/model/response"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
		AgentRegister
	 	Agent 启动时，向 Server 端注册的接口
		req: handle.AgentRegisterReq
*/
func AgentRegister(req request.AgentRegisterReq) (AgentId string, err error) {
	data, err := json.Marshal(req)
	if err != nil {

		return "0", nil
	}
	client := &http.Client{}
	r, _ := http.NewRequest("POST", global.WkgURL+"/v2/agent/register", strings.NewReader(string(data))) // URL-encoded payload
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		return "0", err
	}
	if resp.StatusCode == 200 {
		var res = &response.AgentRegisterRes{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {

			return "0", err
		}
		err = json.Unmarshal(body, res)
		if err != nil {
			return "0", err
		}
		if res.Status == 201 {
			AgentId = res.Data.Id
			fmt.Println("注册成功，探针ID为", res.Data.Id)
			return AgentId, nil
		} else {
			return "0", errors.New("注册失败，失败原因" + res.Msg)
		}
	}
	return "0", errors.New("状态码为" + resp.Status)
}

func ReportUpload(req request.UploadReq) {
	resp, body, errs := POST("/v2/report/upload", req).End()
	if len(errs) > 0 {
		for _, v := range errs {
			fmt.Println(v)
		}
		//fmt.Println("boom")
		return
	}
	if resp.StatusCode == 200 {
		var res response.ResBase
		err := json.Unmarshal([]byte(body), &res)
		if err != nil {
			fmt.Println(err)
			return
		}

		if res.Status == 201 {
			//fmt.Println("pang")
		} else {
			//fmt.Println(res.Msg)
		}

	}
}

func UploadResult(info *request.InputDataStruct) error {
	clien := http.Client{Timeout: time.Second * 5}
	url := global.WkgURL + "/v3/import?token=" + global.V3Token

	body, err := json.Marshal(info)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := clien.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		log.Println("data upload success.")
		Glog.InfoG("data upload success.")
	} else {
		log.Println("data upload failed. respstatus:", resp.Status)
		Glog.InfoG("data upload failed.")
	}
	return nil
}
