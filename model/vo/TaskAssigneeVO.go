package vo

import "github.com/hothotsavage/gstep/model/entity"

type TaskAssigneeVO struct {
	entity.TaskAssignee
	User entity.User `json:"user"`
}
