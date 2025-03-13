package entity

import (
	"github.com/hothotsavage/gstep/util/LocalTime"
	"github.com/hothotsavage/gstep/util/db/entity"
)

type Process struct {
	entity.BaseEntity
	TemplateId  int                  `json:"templateId"`
	StartUserId string               `json:"startUserId"`
	State       string               `json:"state"`
	FinishedAt  *LocalTime.LocalTime `json:"finishedAt"`
}

func (e Process) TableName() string {
	return "process"
}

func (e Process) GetId() any {
	return e.Id
}
