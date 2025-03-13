package vo

import "github.com/hothotsavage/gstep/model/entity"

type ProcessVo struct {
	entity.Process
	Template entity.Template `json:"template"`
}
