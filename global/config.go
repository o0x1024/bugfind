package global

import "bugfind/model/request"

var (
	AttackQueue = make(chan request.UploadReq, 50)
)
