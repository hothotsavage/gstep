package vo

import "github.com/hothotsavage/gstep/model/entity"

type TaskVO struct {
	entity.Task
	Executors []entity.Executor `json:"executors"`
}
