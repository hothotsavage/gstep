package TaskHandler

import (
	"github.com/hothotsavage/gstep/dao/TaskAssigneeDao"
	"github.com/hothotsavage/gstep/dao/TaskDao"
	"github.com/hothotsavage/gstep/enum/TaskState"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/StepService"
	"github.com/hothotsavage/gstep/service/TaskService"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func Pass(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TaskPassDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := DbUtil.GetTx()
	//校验taskid
	pTask := dao.CheckById[entity.Task](dto.TaskId, tx)

	//检查流程提交人是否是候选人
	pProcess := dao.CheckById[entity.Process](pTask.ProcessId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	StepService.CheckCandidate(dto.UserId, &pTemplate.RootStep, pTask.StepId, tx)

	//检查提交人重复提交
	TaskAssigneeDao.CheckAssigneeCanSubmit(pTask.Id, dto.UserId, TaskState.PASS.Code, tx)

	//审核通过
	TaskService.Pass(&dto, tx)

	//任务状态变更通知
	TaskService.NotifyTasksStateChange(pProcess.Id, tx)

	tx.Commit()

	AjaxJson.Success().Response(writer)
}

// 退回到指定上一步
func Retreat(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TaskRefuseDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := DbUtil.GetTx()
	//校验taskid
	pTask := dao.CheckById[entity.Task](dto.TaskId, tx)
	pProcess := dao.CheckById[entity.Process](pTask.ProcessId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)

	//检查提交人重复提交
	TaskAssigneeDao.CheckAssigneeCanSubmit(pTask.Id, dto.UserId, TaskState.PASS.Code, tx)
	//校验提交人在候选人列表中
	StepService.CheckCandidate(dto.UserId, &pTemplate.RootStep, pTask.StepId, tx)
	TaskService.Retreat(&dto, tx)

	//任务状态变更通知
	TaskService.NotifyTasksStateChange(pProcess.Id, tx)

	tx.Commit()

	AjaxJson.Success().Response(writer)
}

func Refuse(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TaskCeaseDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := DbUtil.GetTx()
	//校验taskid
	pTask := dao.CheckById[entity.Task](dto.TaskId, tx)

	//检查流程提交人是否是候选人
	pProcess := dao.CheckById[entity.Process](pTask.ProcessId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	StepService.CheckCandidate(dto.UserId, &pTemplate.RootStep, pTask.StepId, tx)

	//检查提交人重复提交
	TaskAssigneeDao.CheckAssigneeCanSubmit(pTask.Id, dto.UserId, TaskState.REFUSE.Code, tx)

	TaskService.Cease(&dto, tx)

	//任务状态变更通知
	TaskService.NotifyTasksStateChange(pProcess.Id, tx)

	tx.Commit()

	AjaxJson.Success().Response(writer)
}

func Pending(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TaskPendingDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tasks, total := TaskDao.QueryMyPendingTasks(dto.UserId, DbUtil.Db)
	AjaxJson.SuccessByPagination(*tasks, total).Response(writer)
}
