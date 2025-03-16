package TaskDao

import (
	"fmt"
	"github.com/gookit/goutil/strutil"
	"github.com/hothotsavage/gstep/enum/TaskState"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	"gorm.io/gorm"
)

func QueryTaskByStepId(stepId int, processId int, tx *gorm.DB) *entity.Task {
	var detail entity.Task
	tx.Table(detail.TableName()).Where("step_id=? and process_id=?", stepId, processId).First(&detail)
	if 0 == detail.Id {
		return nil
	} else {
		return &detail
	}
}

func QueryMyPendingTasks(userId string, tx *gorm.DB) (*[]entity.Task, int) {
	total := 0
	var details []entity.Task

	err := tx.Raw("select count(1) from task "+
		" where state='started' "+
		" and exists(select 1 from task_assignee ta"+
		" where ta.task_id=task.id "+
		" and ta.user_id=?)", userId).Scan(&total).Error
	if nil != err {
		msg := fmt.Sprintf("找不到待处理任务: %s", err)
		panic(ServerError.New(msg))
	}
	err = tx.Raw("select * from task "+
		" where state='started' "+
		" and exists(select 1 from task_assignee ta"+
		" where ta.task_id=task.id "+
		" and ta.user_id=?)", userId).Scan(&details).Error
	if nil != err {
		msg := fmt.Sprintf("找不到待处理任务: %s", err)
		panic(ServerError.New(msg))
	}
	return &details, total
}

// 查询流程实例的任务列表
func Query(taskQueryDto dto.TaskQueryDto, tx *gorm.DB) []entity.Task {
	var tasks []entity.Task
	sql := "select * from task where 1=1 "
	if taskQueryDto.ProcessId > 0 {
		sql += fmt.Sprintf(" and process_id=%d ", taskQueryDto.ProcessId)
	}
	if taskQueryDto.StartStepId > 0 {
		sql += fmt.Sprintf(" and step_id>=%d ", taskQueryDto.StartStepId)
	}
	if strutil.IsNotBlank(taskQueryDto.State) {
		sql += fmt.Sprintf(" and state='%s' ", taskQueryDto.State)
	}
	err := tx.Raw(sql).Scan(&tasks).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程(processId=%d)任务失败: %s", taskQueryDto.ProcessId, err)
		panic(ServerError.New(msg))
	}
	return tasks
}

func GetStartedTask(processId int, tx *gorm.DB) entity.Task {
	var task entity.Task
	err := tx.Raw("select * from task where process_id=? and state=? order by step_id asc", processId, TaskState.STARTED.Code).First(&task).Error
	if nil != err {
		msg := fmt.Sprintf("查询待审核任务(processId=%d)失败: %s", processId, err)
		panic(ServerError.New(msg))
	}
	return task
}

// 删除指定步骤之后的未开始任务
func DeleteNextUnstartTasks(processId int, startStepId int, tx *gorm.DB) {
	err := tx.Exec("delete from task where process_id=? and step_id>=? and state=?", processId, startStepId, TaskState.UNSTART.Code).Error
	if nil != err {
		msg := fmt.Sprintf("删除下一步任务列表失败(processId=%d,stepId=%d)失败: %s", processId, startStepId, err)
		panic(ServerError.New(msg))
	}
}
