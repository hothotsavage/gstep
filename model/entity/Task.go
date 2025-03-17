package entity

import "github.com/hothotsavage/gstep/util/db/entity"

type Task struct {
	entity.BaseEntity
	ProcessId   int
	Form        *map[string]any `json:"form" gorm:"serializer:json"`
	AuditMethod string
	StepId      int
	Title       string
	Category    string
	State       string
	Candidates  []string `json:"candidates" gorm:"serializer:json"`
}

func (e Task) TableName() string {
	return "task"
}

func (e Task) GetId() any {
	return e.Id
}
