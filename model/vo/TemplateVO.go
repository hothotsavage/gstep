package vo

import "github.com/hothotsavage/gstep/model/entity"

type TemplateVO struct {
	entity.Template
	ProcessCount int `json:"processCount"`
}
