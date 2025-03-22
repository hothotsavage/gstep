package entity

import (
	"github.com/hothotsavage/gstep/util/db/entity"
)

type Mould struct {
	entity.BaseEntity
	Title string `json:"title"`
}

func (e Mould) TableName() string {
	return "mould"
}

func (e Mould) GetId() any {
	return e.Id
}
