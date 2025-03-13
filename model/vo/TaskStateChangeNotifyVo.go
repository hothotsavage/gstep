package vo

import "github.com/hothotsavage/gstep/model/entity"

type TaskStateChangeNotifyVo struct {
	Process entity.Process                `json:"process"`
	Tasks   []TaskStateChangeNotifyTaskVo `json:"tasks"`
}
