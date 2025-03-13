package vo

import "github.com/hothotsavage/gstep/model/entity"

type TaskStateChangeNotifyTaskVo struct {
	Task      entity.Task           `json:"task"`
	Assignees []entity.TaskAssignee `json:"assignees"`
}
