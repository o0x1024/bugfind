package global

import "BugFind/model/request"

var (
	AttackQueue = make(chan request.UploadReq, 50)
)
