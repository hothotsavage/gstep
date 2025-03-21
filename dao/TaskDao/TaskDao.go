package TaskDao

import (
	"fmt"
	"github.com/gookit/goutil/strutil"
	"github.com/hothotsavage/gstep/enum/StepCat"
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
	if taskQueryDto.StartTaskId > 0 {
		sql += fmt.Sprintf(" and id>=%d ", taskQueryDto.StartTaskId)
	}
	if strutil.IsNotBlank(taskQueryDto.State) {
		sql += fmt.Sprintf(" and state='%s' ", taskQueryDto.State)
	}
	if strutil.IsNotBlank(taskQueryDto.Category) {
		sql += fmt.Sprintf(" and category='%s' ", taskQueryDto.Category)
	}
	sql += fmt.Sprintf(" order by id asc")
	err := tx.Raw(sql).Scan(&tasks).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程(processId=%d)任务失败: %s", taskQueryDto.ProcessId, err)
		panic(ServerError.New(msg))
	}
	return tasks
}

func GetStartedTask(processId int, tx *gorm.DB) entity.Task {
	var task entity.Task
	err := tx.Raw("select * from task where process_id=? and state=? order by id asc", processId, TaskState.STARTED.Code).First(&task).Error
	if nil != err {
		msg := fmt.Sprintf("查询待审核任务(processId=%d)失败: %s", processId, err)
		panic(ServerError.New(msg))
	}
	return task
}

func GetFirstTask(processId int, tx *gorm.DB) entity.Task {
	var task entity.Task
	err := tx.Raw("select * from task where process_id=? order by id asc limit 1", processId).First(&task).Error
	if nil != err {
		msg := fmt.Sprintf("查询第一个任务(processId=%d)失败: %s", processId, err)
		panic(ServerError.New(msg))
	}
	return task
}

// 查询流程实例的所有任务都已审核结束
func IsProcessFinish(processId int, tx *gorm.DB) bool {
	//查询没有未结束的任务，则流程实例结束
	cnt := 0
	err := tx.Raw("select count(1) from task "+
		" where process_id = ? and (state=? or state=?)", processId, TaskState.STARTED.Code, TaskState.UNSTART.Code).Scan(&cnt).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程实例的started任务数量失败: %s", err)
		panic(ServerError.New(msg))
	}
	return cnt < 1
}

// 删除未开始的任务列表和执行人列表
func DeleteUnstartTasksAndExecutors(processId int, tx *gorm.DB) {
	err := tx.Exec("delete from executor "+
		" where 1=1 "+
		" and process_id=? "+
		" and exists(select 1 "+
		" from task t "+
		" where t.id=executor.task_id "+
		" and t.state=?) ", processId, TaskState.UNSTART.Code).Error
	if nil != err {
		msg := fmt.Sprintf("删除未启动任务的执行人列表失败(processId=%d)失败: %s", processId, err)
		panic(ServerError.New(msg))
	}

	err = tx.Exec("delete from task where process_id=? and state=?", processId, TaskState.UNSTART.Code).Error
	if nil != err {
		msg := fmt.Sprintf("删除未启动任务列表失败(processId=%d)失败: %s", processId, err)
		panic(ServerError.New(msg))
	}
}

// 查询可回退的步骤id列表
func GetRefusePrevSteps(processId int, tx *gorm.DB) []int {
	var ids []int
	maxRefuseId := 0
	err := tx.Raw("select ifnull(max(id),0) from task "+
		" where process_id = ? and state=?", processId, TaskState.REFUSE.Code).Scan(&maxRefuseId).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程(processId=%d)的最大拒绝taskId失败: %s", processId, err)
		panic(ServerError.New(msg))
	}
	err = tx.Raw("select distinct step_id from task where process_id=? "+
		" and state=?"+
		" and (category=? or category=?) "+
		" and id>?", processId, TaskState.PASS.Code, StepCat.AUDIT.Code, StepCat.START.Code, maxRefuseId).Scan(&ids).Error
	if nil != err {
		msg := fmt.Sprintf("查询可回退的步骤id列表(processId=%d)失败: %s", processId, err)
		panic(ServerError.New(msg))
	}
	return ids
}

// 查询流程的回退步骤的任务
func GetPrevTaskByStepId(processId int, prevStepId int, tx *gorm.DB) entity.Task {
	tasks := []entity.Task{}
	err := tx.Raw("select * from task where process_id=? "+
		" and state=?"+
		" and step_id=? ", processId, TaskState.PASS.Code, prevStepId).Scan(&tasks).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程(processId=%d)的回退步骤(stepId=%d)的任务列表失败: %s", processId, prevStepId, err)
		panic(ServerError.New(msg))
	}
	if len(tasks) == 0 {
		msg := fmt.Sprintf("找不到流程(processId=%d)的回退步骤(stepId=%d)的任务列表", processId, prevStepId)
		panic(ServerError.New(msg))
	}

	return tasks[0]
}
