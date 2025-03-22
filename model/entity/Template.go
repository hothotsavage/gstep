package entity

import (
	"github.com/hothotsavage/gstep/util/db/entity"
)

type Template struct {
	entity.BaseEntity
	MouldId  int     `json:"mouldId"`
	Title    string  `json:"title"`
	Version  int     `json:"version"`
	RootStep Step    `json:"rootStep" gorm:"serializer:json"`
	Fields   []Field `json:"fields" gorm:"serializer:json"`
	State    string  `json:"state"`
}

func (e Template) TableName() string {
	return "template"
}

func (e Template) GetId() any {
	return e.Id
}
