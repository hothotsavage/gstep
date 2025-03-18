package vo

import "github.com/hothotsavage/gstep/model/entity"

type TaskVO struct {
	entity.Task
	CandidateUsers []entity.User    `json:"candidateUsers"`
	Assignees      []TaskAssigneeVO `json:"assignees"`
}
