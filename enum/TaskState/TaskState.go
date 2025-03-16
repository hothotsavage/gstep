package TaskState

import "github.com/hothotsavage/gstep/util/enum"

type TaskState struct {
	enum.BaseEnum[string]
}

var UNSTART = TaskState{}
var STARTED = TaskState{}
var PASS = TaskState{}
var REFUSE = TaskState{}

func init() {
	UNSTART.Code = "unstart"
	UNSTART.Title = "未开始"

	STARTED.Code = "started"
	STARTED.Title = "开始"

	PASS.Code = "pass"
	PASS.Title = "同意"

	REFUSE.Code = "refuse"
	REFUSE.Title = "拒绝"

}
