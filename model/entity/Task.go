package entity

import "github.com/hothotsavage/gstep/util/db/entity"

type Task struct {
	entity.BaseEntity
	ProcessId   int             `json:"processId"`
	Form        *map[string]any `json:"form" gorm:"serializer:json"`
	AuditMethod string          `json:"auditMethod"`
	StepId      int             `json:"stepId"`
	Title       string          `json:"title"`
	Category    string          `json:"category"`
	State       string          `json:"state"`
}

func (e Task) TableName() string {
	return "task"
}

func (e Task) GetId() any {
	return e.Id
}
