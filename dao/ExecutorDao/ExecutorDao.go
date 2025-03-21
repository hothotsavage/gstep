package ExecutorDao

import (
	"fmt"
	"github.com/hothotsavage/gstep/enum/TaskState"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	"gorm.io/gorm"
)

func GetTaskExecutors(taskId int, tx *gorm.DB) []entity.Executor {
	var executors []entity.Executor
	sql := "select * from executor where 1=1 " +
		" and task_id=? " +
		" order by id asc "
	err := tx.Raw(sql, taskId).Scan(&executors).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程任务执行人列表(taskId=%d)失败: %s", taskId, err)
		panic(ServerError.New(msg))
	}
	return executors
}

func GetTaskExecutor(taskId int, userId string, tx *gorm.DB) *entity.Executor {
	var executors []entity.Executor
	sql := "select * from executor where 1=1 " +
		" and task_id=? " +
		" and user_id=? " +
		" order by id asc "
	err := tx.Raw(sql, taskId, userId).Scan(&executors).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程任务执行人(taskId=%d)失败: %s", taskId, err)
		panic(ServerError.New(msg))
	}
	if len(executors) < 1 {
		return nil
	}
	return &executors[0]
}

func PassCount(taskId int, tx *gorm.DB) int {
	var count int64
	tx.Table("executor").Where("task_id=? and state=?", taskId, TaskState.PASS.Code).Count(&count)
	return int(count)
}

func ExecutorCount(taskId int, tx *gorm.DB) int {
	var count int64
	tx.Table("executor").Where("task_id=?", taskId).Count(&count)
	return int(count)
}

func DeleteTaskExecutors(taskId int, tx *gorm.DB) {
	err := tx.Exec("delete from executor where task_id=?", taskId).Error
	if nil != err {
		msg := fmt.Sprintf("删除任务执行人列表失败(taskId=%d)失败: %s", taskId, err)
		panic(ServerError.New(msg))
	}
}
