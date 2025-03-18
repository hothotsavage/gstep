package vo

import "github.com/hothotsavage/gstep/model/entity"

type ProcessVO struct {
	entity.Process
	Template entity.Template `json:"template"`
}
