package ProcessService

import (
	"fmt"
	"github.com/hothotsavage/gstep/dao/TaskAssigneeDao"
	"github.com/hothotsavage/gstep/dao/TaskDao"
	"github.com/hothotsavage/gstep/dao/TemplateDao"
	"github.com/hothotsavage/gstep/enum/AuditMethodCat"
	"github.com/hothotsavage/gstep/enum/ProcessState"
	"github.com/hothotsavage/gstep/enum/StepCat"
	"github.com/hothotsavage/gstep/enum/TaskState"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/model/vo"
	"github.com/hothotsavage/gstep/service/StepService"
	"github.com/hothotsavage/gstep/service/TaskService"
	"github.com/hothotsavage/gstep/util/ExpressionUtil"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func Start(processStartDto *dto.ProcessStartDto, tx *gorm.DB) vo.NotifyVO {
	process := entity.Process{}
	copier.Copy(process, processStartDto)

	//创建流程
	pTemplate := TemplateDao.GetLatestVersion(processStartDto.MouldId, tx)
	if nil == pTemplate {
		panic(ServerError.New("无效的模板"))
	}

	//校验流程发起人是否在候选人列表中
	StepService.CheckStepCandidate(processStartDto.UserId, processStartDto.Form, pTemplate.Id, 1, tx)

	process.TemplateId = pTemplate.Id
	process.StartUserId = processStartDto.UserId
	process.State = ProcessState.STARTED.Code
	dao.SaveOrUpdate(&process, tx)
	processVO := ToVO(&process, tx)

	//创建启动任务
	pStartTask, pTaskAssignee := TaskService.NewStartTask(&process, processStartDto.UserId, processStartDto.Form, processStartDto.Memo, tx)

	notifyUserIds := []string{}
	//创建后续的任务列表
	nextStep := processVO.Template.RootStep.NextStep
	if nil != nextStep && nextStep.Category != StepCat.END.Code {
		notifyUserIds = TaskService.MakeTasks(process.Id, nextStep.Id, processStartDto.Form, tx)
	}

	//任务状态变更通知
	//TaskService.NotifyTasksStateChange(process.Id)

	//生成通知消息文案
	notifyMessage := TaskService.MakeNotifyMessage(*pTaskAssignee, tx)

	notifyVO := vo.NotifyVO{
		process.Id,
		pStartTask.Id,
		notifyUserIds,
		notifyMessage,
	}
	return notifyVO
}

// 审核通过
func Pass(processPassDto dto.ProcessPassDto, tx *gorm.DB) vo.NotifyVO {
	//查询已启动的任务
	startedTask := TaskDao.GetStartedTask(processPassDto.ProcessId, tx)
	//校验taskid
	pStartedTask := dao.CheckById[entity.Task](startedTask.Id, tx)
	pProcess := dao.CheckById[entity.Process](pStartedTask.ProcessId, tx)
	startedStep := GetStep(processPassDto.ProcessId, pStartedTask.StepId, tx)

	if startedStep.Id < 1 {
		panic(ServerError.New("无效的流程步骤"))
	}
	if startedStep.Category == StepCat.END.Code {
		panic(ServerError.New("结束步骤不用提交"))
	}

	//更新任务候选人
	pStartedTask.Candidates = TaskService.MakeCandidates(startedStep, processPassDto.Form, tx)
	dao.SaveOrUpdate(pStartedTask, tx)
	//检查候选人
	TaskService.CheckCandidate(processPassDto.UserId, processPassDto.Form, pStartedTask.Id, tx)
	//检查提交人重复提交
	TaskAssigneeDao.CheckAssigneeCanSubmit(pStartedTask.Id, processPassDto.UserId, tx)

	//保存任务提交人
	submitIndex := TaskAssigneeDao.GetMaxSubmitIndex(pStartedTask.Id, tx) + 1
	assignee := entity.TaskAssignee{}
	assignee.ProcessId = pStartedTask.ProcessId
	assignee.StepId = pStartedTask.StepId
	assignee.TaskId = pStartedTask.Id
	assignee.UserId = processPassDto.UserId
	assignee.State = TaskState.PASS.Code
	assignee.Form = processPassDto.Form
	assignee.SubmitIndex = submitIndex
	assignee.Memo = processPassDto.Memo
	dao.SaveOrUpdate(&assignee, tx)

	//保存任务表单
	pStartedTask.Form = processPassDto.Form
	//更新任务状态
	if CanTaskPass(pStartedTask, pProcess, tx) {
		pStartedTask.State = TaskState.PASS.Code
	}
	dao.SaveOrUpdate(pStartedTask, tx)

	notifyUserIds := []string{}
	//创建后续的任务列表
	nextStep := startedStep.NextStep
	if nil != nextStep && nextStep.Id > 0 && nextStep.Category != StepCat.END.Code {
		//先删除当前步骤之后的所有未开始的任务
		TaskDao.DeleteUnstartTasks(pStartedTask.ProcessId, tx)
		//创建后续任务
		notifyUserIds = TaskService.MakeTasks(pStartedTask.ProcessId, nextStep.Id, processPassDto.Form, tx)
	}

	//更新流程实例状态
	UpdateProcessState(processPassDto.ProcessId, tx)

	//任务状态变更通知
	//TaskService.NotifyTasksStateChange(pProcess.Id)

	//生成通知消息文案
	notifyMessage := TaskService.MakeNotifyMessage(assignee, tx)

	notifyVO := vo.NotifyVO{
		pProcess.Id,
		pStartedTask.Id,
		notifyUserIds,
		notifyMessage,
	}
	return notifyVO
}

// 拒绝
func Refuse(processRefuseDto dto.ProcessRefuseDto, tx *gorm.DB) vo.NotifyVO {
	//查询已启动的任务
	startedTask := TaskDao.GetStartedTask(processRefuseDto.ProcessId, tx)
	//校验taskid
	pStartedTask := dao.CheckById[entity.Task](startedTask.Id, tx)
	dao.CheckById[entity.Process](pStartedTask.ProcessId, tx)
	startedStep := GetStep(processRefuseDto.ProcessId, pStartedTask.StepId, tx)
	GetNextStep(processRefuseDto.ProcessId, processRefuseDto.PrevStepId, tx)

	if startedStep.Id < 1 {
		panic(ServerError.New("无效的流程步骤"))
	}
	if startedStep.Category == StepCat.END.Code {
		panic(ServerError.New("结束步骤不用提交"))
	}

	//更新任务候选人
	pStartedTask.Candidates = TaskService.MakeCandidates(startedStep, processRefuseDto.Form, tx)
	dao.SaveOrUpdate(pStartedTask, tx)
	//检查候选人
	TaskService.CheckCandidate(processRefuseDto.UserId, processRefuseDto.Form, pStartedTask.Id, tx)
	//检查提交人重复提交
	TaskAssigneeDao.CheckAssigneeCanSubmit(pStartedTask.Id, processRefuseDto.UserId, tx)

	//保存任务提交人
	submitIndex := TaskAssigneeDao.GetMaxSubmitIndex(pStartedTask.Id, tx) + 1
	assignee := entity.TaskAssignee{}
	assignee.ProcessId = pStartedTask.ProcessId
	assignee.StepId = pStartedTask.StepId
	assignee.TaskId = pStartedTask.Id
	assignee.State = TaskState.REFUSE.Code
	assignee.Form = processRefuseDto.Form
	assignee.SubmitIndex = submitIndex
	assignee.UserId = processRefuseDto.UserId
	assignee.Memo = processRefuseDto.Memo
	dao.SaveOrUpdate(&assignee, tx)

	//保存任务表单
	pStartedTask.Form = processRefuseDto.Form
	//更新任务状态
	pStartedTask.State = TaskState.REFUSE.Code
	dao.SaveOrUpdate(pStartedTask, tx)

	notifyUserIds := []string{}
	//创建后续的任务列表
	nextStep := startedStep.NextStep
	if nil != nextStep && nextStep.Id > 0 && nextStep.Category != StepCat.END.Code {
		//先删除当前步骤之后的所有未开始的任务
		TaskDao.DeleteUnstartTasks(pStartedTask.ProcessId, tx)
	}

	//创建回退步骤及后续任务
	notifyUserIds = TaskService.MakeTasks(pStartedTask.ProcessId, processRefuseDto.PrevStepId, processRefuseDto.Form, tx)

	//更新流程实例状态
	UpdateProcessState(processRefuseDto.ProcessId, tx)

	//任务状态变更通知
	//TaskService.NotifyTasksStateChange(pProcess.Id)

	//生成通知消息文案
	notifyMessage := TaskService.MakeNotifyMessage(assignee, tx)

	notifyVO := vo.NotifyVO{
		processRefuseDto.ProcessId,
		pStartedTask.Id,
		notifyUserIds,
		notifyMessage,
	}
	return notifyVO
}

func ToVO(pProcess *entity.Process, tx *gorm.DB) vo.ProcessVO {
	aVo := vo.ProcessVO{}
	if nil == pProcess {
		return aVo
	}

	aVo.Process = *pProcess
	template := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	aVo.Template = *template

	return aVo
}

func ToDetailVO(pProcess *entity.Process, tx *gorm.DB) vo.ProcessDetailVO {
	aVo := vo.ProcessDetailVO{}
	if nil == pProcess {
		return aVo
	}

	aVo.Process = *pProcess
	template := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	aVo.Template = *template

	tasks := TaskService.GetTasksByProcessId(pProcess.Id, tx)
	taskVOs := []vo.TaskVO{}
	for _, task := range tasks {
		taskVO := TaskService.ToVO(task, tx)
		taskVOs = append(taskVOs, taskVO)
	}
	aVo.Tasks = taskVOs
	return aVo
}

func GetStep(processId int, stepId int, tx *gorm.DB) entity.Step {
	pPrcess := dao.CheckById[entity.Process](processId, tx)
	processVO := ToVO(pPrcess, tx)
	pStep := StepService.FindStep(&processVO.Template.RootStep, stepId, tx)

	if nil == pStep {
		panic(ServerError.New(fmt.Sprintf("无效的步骤id: %d", stepId)))
	}

	return *pStep
}

func GetSteps(processId int, stepIds []int, tx *gorm.DB) []entity.Step {
	pPrcess := dao.CheckById[entity.Process](processId, tx)
	processVO := ToVO(pPrcess, tx)
	steps := []entity.Step{}
	for _, stepId := range stepIds {
		pStep := StepService.FindStep(&processVO.Template.RootStep, stepId, tx)
		if nil == pStep {
			panic(ServerError.New(fmt.Sprintf("无效的步骤id: %d", stepId)))
		}
		steps = append(steps, *pStep)
	}

	return steps
}

func GetNextStep(processId int, startStepId int, tx *gorm.DB) entity.Step {
	pPrcess := dao.CheckById[entity.Process](processId, tx)
	processVO := ToVO(pPrcess, tx)
	pStep := StepService.FindStep(&processVO.Template.RootStep, startStepId, tx)
	if nil == pStep || pStep.Id < 1 {
		panic(ServerError.New("无效的流程步骤id"))
	}

	return *pStep.NextStep
}

// 判断是否可通过
func CanTaskPass(pTask *entity.Task, pProcess *entity.Process, tx *gorm.DB) bool {
	if pTask.Category == StepCat.CONDITION.Code {
		pStep := GetStep(pProcess.Id, pTask.StepId, tx)
		exp := ExpressionUtil.Template2jsExpression(pStep.Expression, pTask.Form)
		isPass := ExpressionUtil.RunJsExpression(exp)
		return isPass
	} else {
		passCount := TaskAssigneeDao.PassCount(pTask.Id, tx)

		//或签
		if pTask.AuditMethod == AuditMethodCat.OR.Code {
			return passCount > 0
		} else { //会签
			candidateCount := StepService.CandidateCount(pTask.Id, tx)
			return passCount >= candidateCount
		}
	}
}

func UpdateProcessState(processId int, tx *gorm.DB) {
	pProcess := dao.CheckById[entity.Process](processId, tx)

	if IsFinish(processId, tx) {
		pProcess.State = ProcessState.FINISH_PASS.Code
		dao.SaveOrUpdate(pProcess, tx)
	}
}

func IsFinish(processId int, tx *gorm.DB) bool {
	isFinish := TaskDao.IsProcessFinish(processId, tx)
	return isFinish
}

// 查询可回退的步骤列表
func RefusePrevSteps(processId int, tx *gorm.DB) []entity.Step {
	stepIds := TaskDao.GetRefusePrevSteps(processId, tx)
	steps := GetSteps(processId, stepIds, tx)
	return steps
}
