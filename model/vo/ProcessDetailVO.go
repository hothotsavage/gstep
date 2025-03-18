package vo

import "github.com/hothotsavage/gstep/model/entity"

type ProcessDetailVO struct {
	entity.Process
	Template entity.Template `json:"template"`
	Tasks    []TaskVO        `json:"tasks"`
}
