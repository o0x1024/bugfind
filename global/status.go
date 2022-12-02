package global

import (
	types2 "BugFind/model/types"
)

var (
	Target        = make(chan string, 500)
	WkgURL        = "http://42.193.253.66:7788"
	NewDomainPath = "./result.txt"

	Alarm    = make(chan types2.Alarm, 1000)
	FofaDone = make(chan bool)
	Version  = "v1.0"
	AgentId  = ""
	V3Token  = "123123123123"
)
