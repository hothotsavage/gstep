package entity

import "github.com/hothotsavage/gstep/util/db/entity"

type Executor struct {
	entity.BaseEntity
	ProcessId   int             `json:"processId"`
	StepId      int             `json:"stepId"`
	TaskId      int             `json:"taskId"`
	UserId      string          `json:"userId"`
	State       string          `json:"state"`
	SubmitIndex int             `json:"submitIndex"`
	Form        *map[string]any `json:"form" gorm:"serializer:json"`
	Memo        string          `json:"memo"`
}

func (e Executor) TableName() string {
	return "executor"
}

func (e Executor) GetId() any {
	return e.Id
}
