package core

import (
	"bugfind/global"
	"bugfind/internal/handle"
	"bugfind/internal/libs/safe"
	"bugfind/internal/service"
	"bugfind/internal/utils"
	"bugfind/model/request"
	"encoding/json"
	"fmt"

	"log"
	"strconv"
	"sync"
	"time"
)

var agentInstance agentInstanceType

type agentInstanceType struct {
	startOnce          sync.Once
	instanceAccessLock sync.RWMutex
}

func Start() {
	agentInstance.start()
}

func (p *agentInstanceType) start() {
	fmt.Println("[+] agent start...")
	var err error

	p.startOnce.Do(func() {
		safe.Go(func() error {
			//agent注册
			err = service.AgentRegister()
			if err != nil {
				log.Fatalln("agent register failed.")
				return err
			}
			//hook初始化

			for {
				err := safe.Call(func() error {
					err = safe.Call(p.Serve)
					if err == nil {
						return nil
					}
					return err
				})
				if err == nil {
					return nil
				}
			}
		}, nil)
	})
}

func (ag *agentInstanceType) Serve() error {

	ticker := time.NewTicker(time.Duration(3) * time.Minute)

	for {
		select {
		case <-ticker.C:
			HeartBeat()
			//case attackReq := <-global.AttackQueue:
			//检查到攻击和异常信息上报服务器
			//handle.ReportUpload(attackReq)
			//if global.Config.LogStatus && logOpenStatus {
			//	raspOut.WriteString( attackReq.String() + "\n")
			//}
			//default:

		}
	}
	return nil
}

func HeartBeat() {
	s, err := utils.GetServerInfo()
	if err != nil {
		return
	}
	var req request.UploadReq

	cpuMap := make(map[string]string)
	memoryMap := make(map[string]string)
	var cpus float64 = 0
	for _, v := range s.Cpu.Cpus {
		cpus += v
	}
	cpuRate := fmt.Sprintf("%.2f", cpus/float64(len(s.Cpu.Cpus)))
	memoryRate := fmt.Sprintf("%.2f", float64(s.Rrm.UsedMB)/float64(s.Rrm.TotalMB))
	total := strconv.Itoa(s.Rrm.TotalMB) + "MB"
	use := strconv.Itoa(s.Rrm.UsedMB) + "MB"
	cpuMap["rate"] = cpuRate
	memoryMap["total"] = total
	memoryMap["use"] = use
	memoryMap["rate"] = memoryRate
	cpuByte, _ := json.Marshal(cpuMap)
	memoryByte, _ := json.Marshal(memoryMap)

	req.Type = 1
	req.Detail.Pant.Disk = "{}"
	req.Detail.Pant.Cpu = string(cpuByte)
	req.Detail.Pant.Memory = string(memoryByte)
	req.Detail.AgentId = global.AgentId
	handle.ReportUpload(req)
}
