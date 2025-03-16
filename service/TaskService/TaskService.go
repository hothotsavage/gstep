package TaskService

import (
	"fmt"
	"github.com/hothotsavage/gstep/config"
	"github.com/hothotsavage/gstep/dao/TaskAssigneeDao"
	"github.com/hothotsavage/gstep/dao/TaskDao"
	"github.com/hothotsavage/gstep/dao/UserDao"
	"github.com/hothotsavage/gstep/enum/AuditMethodCat"
	"github.com/hothotsavage/gstep/enum/CandidateCat"
	"github.com/hothotsavage/gstep/enum/ProcessState"
	"github.com/hothotsavage/gstep/enum/StepCat"
	"github.com/hothotsavage/gstep/enum/TaskState"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/model/vo"
	"github.com/hothotsavage/gstep/service/StepService"
	"github.com/hothotsavage/gstep/util/CONSTANT"
	"github.com/hothotsavage/gstep/util/CollectionUtil"
	"github.com/hothotsavage/gstep/util/ExpressionUtil"
	"github.com/hothotsavage/gstep/util/JsonUtil"
	"github.com/hothotsavage/gstep/util/LocalTime"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestUtil"
	"gorm.io/gorm"
	"log"
	"time"
)

// 终止流程
func Cease(pDto *dto.TaskCeaseDto, tx *gorm.DB) int {
	pTask := dao.CheckById[entity.Task](pDto.TaskId, tx)
	pProcess := dao.CheckById[entity.Process](pTask.ProcessId, tx)

	//保存任务提交人
	submitIndex := TaskAssigneeDao.GetMaxSubmitIndex(pTask.Id, tx) + 1
	assignee := entity.TaskAssignee{}
	assignee.TaskId = pTask.Id
	assignee.UserId = pDto.UserId
	assignee.State = TaskState.REFUSE.Code
	assignee.SubmitIndex = submitIndex
	dao.SaveOrUpdate(&assignee, tx)

	//保存任务表单
	pTask.Form = pDto.Form
	//更新任务状态
	pTask.State = TaskState.REFUSE.Code
	dao.SaveOrUpdate(pTask, tx)

	//撤销上一步
	pProcess.State = ProcessState.FINISH_REFUSE.Code
	dao.SaveOrUpdate(pProcess, tx)

	return submitIndex
}

// 查找流程的下一个步骤
// 碰到分支步骤,取满足条件的条件步骤
// 碰到条件步骤,取条件步骤的下一个步骤
func GetNextStep(currentStepId int, pTemplate *entity.Template, form *map[string]any, tx *gorm.DB) *entity.Step {
	pStep := StepService.FindStep(&pTemplate.RootStep, currentStepId)

	if nil == pStep {
		panic(ServerError.New("找不到流程步骤"))
	}

	//分支步骤,找满足条件的子条件步骤的下一步
	if pStep.Category == StepCat.BRANCH.Code {
		if len(pStep.BranchSteps) < 2 {
			panic(ServerError.New("分支步骤的分支数量小于2"))
		}

		var pDefaultConditionStep = &entity.Step{}
		//先找满足的非默认条件
		for _, v := range pStep.BranchSteps {
			if nil == v {
				panic(ServerError.New("无效流程步骤"))
			}

			if v.Category != StepCat.CONDITION.Code {
				panic(ServerError.New("流程分支的首个步骤不是条件类型步骤"))
			}

			if v.Title == "默认条件" {
				pDefaultConditionStep = v
				continue
			}

			isPass := ExpressionUtil.ExecuteExpression(v.Expression, form)
			if isPass {
				nextStep := GetNextStep(v.Id, pTemplate, form, tx)
				return nextStep
			}
		}
		//所有条件都不满足,走默认条件步骤
		nextStep := GetNextStep(pDefaultConditionStep.Id, pTemplate, form, tx)
		return nextStep
	}

	//没有下一个步骤,往前取分支步骤的下一步
	if nil != pStep.NextStep && pStep.NextStep.Id != 0 {
		//非条件步骤,返回下一步
		if !StepCat.IsRoute(pStep.NextStep.Category) {
			return pStep.NextStep
		}

		pNextStep := GetNextStep(pStep.NextStep.Id, pTemplate, form, tx)
		return pNextStep
	} else {
		//从父分支步骤开始往前递归查找有下一步的分支步骤
		pPrevBranchStep := StepService.FindPrevBranchStepWithNextStep(&pTemplate.RootStep, pStep.Id)
		if nil != pPrevBranchStep.NextStep {
			return pPrevBranchStep.NextStep
		}
	}

	return nil
}

func NewTaskByStep(pStep *entity.Step, pProcess *entity.Process, submitIndex int, form *map[string]any, tx *gorm.DB) *entity.Task {
	if nil == pStep {
		panic(ServerError.New("流程步骤不能为空"))
	}

	if pStep.Category == StepCat.BRANCH.Code {
		panic(ServerError.New("无法用分支步骤创建流程任务"))
	} else if pStep.Category == StepCat.CONDITION.Code {
		panic(ServerError.New("无法用条件步骤创建流程任务"))
	} else if pStep.Category == StepCat.END.Code {
		panic(ServerError.New("无法用结束步骤创建流程任务"))
	}

	//if len(pStep.Candidates) == 0 {
	//	panic(ServerError.New("不能用无候选人的流程步骤创建流程任务"))
	//}

	//若没有退回的任务,则创建新任务
	pTask := TaskDao.QueryTaskByStepId(pStep.Id, pProcess.Id, tx)
	if nil == pTask {
		pTask = &entity.Task{}
		pTask.ProcessId = pProcess.Id
		pTask.StepId = pStep.Id
		pTask.Title = pStep.Title
		pTask.Category = pStep.Category
		pTask.AuditMethod = pStep.AuditMethod
		pTask.Form = form
	}
	pTask.State = TaskState.UNSTART.Code
	dao.SaveOrUpdate(pTask, tx)

	//创建任务的候选人
	if pStep.Category != StepCat.END.Code {
		for _, c := range pStep.Candidates {
			if c.Category == CandidateCat.USER.Code {
				assignee := entity.TaskAssignee{}
				assignee.TaskId = pTask.Id
				assignee.UserId = c.Value
				assignee.State = pTask.State
				assignee.SubmitIndex = submitIndex
				assignee.Form = form
				dao.SaveOrUpdate(&assignee, tx)
			} else if c.Category == CandidateCat.DEPARTMENT.Code {
				users := UserDao.GetGrandsonDepartmentUsers(c.Value, tx)
				for _, user := range users {
					assignee := entity.TaskAssignee{}
					assignee.TaskId = pTask.Id
					assignee.UserId = user.Id
					assignee.State = pTask.State
					assignee.SubmitIndex = submitIndex
					assignee.Form = form
					dao.SaveOrUpdate(&assignee, tx)
				}
			}
		}
	}

	return pTask
}

// 创建启动任务
func NewStartTask(pProcess *entity.Process, startUserId string, form *map[string]any, tx *gorm.DB) *entity.Task {
	//创建启动任务
	task := entity.Task{}
	task.ProcessId = pProcess.Id

	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	rootStep := pTemplate.RootStep

	//检查流程提交人是否是候选人
	StepService.CheckCandidate(startUserId, form, &rootStep, task.StepId, tx)

	//创建启动任务
	task.StepId = rootStep.Id
	task.Title = rootStep.Title
	task.Form = form
	task.Category = rootStep.Category
	task.State = TaskState.PASS.Code
	if len(rootStep.AuditMethod) == 0 {
		task.AuditMethod = AuditMethodCat.OR.Code
	} else {
		task.AuditMethod = rootStep.AuditMethod
	}
	dao.SaveOrUpdate(&task, tx)

	//创建启动任务的候选人
	assignee := entity.TaskAssignee{}
	assignee.TaskId = task.Id
	assignee.UserId = startUserId
	assignee.State = TaskState.PASS.Code
	assignee.SubmitIndex = 1
	assignee.Form = form
	dao.SaveOrUpdate(&assignee, tx)

	return &task
}

// 审核通过流程
func FinishPassProcess(pProcess *entity.Process, tx *gorm.DB) {
	pProcess.State = ProcessState.FINISH_PASS.Code
	finishTime := LocalTime.LocalTime(time.Now())
	pProcess.FinishedAt = &finishTime
	dao.SaveOrUpdate(pProcess, tx)
}

// 创建指定步骤之后的所有任务列表
func MakeTasks(processId int, startStepId int, form *map[string]any, tx *gorm.DB) {
	pProcess := dao.CheckById[entity.Process](processId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	pStartStep := StepService.FindStep(&pTemplate.RootStep, startStepId)
	NewTaskByStep(pStartStep, pProcess, 1, nil, tx)
	//根据form表单，创建开始步骤之后的任务列表
	stepId := pStartStep.Id
	for {
		pNextStep := GetNextStep(stepId, pTemplate, form, tx)
		if nil == pNextStep || 0 == pNextStep.Id {
			panic(ServerError.New(fmt.Sprintf("找不到步骤(stepId=%s)下一个步骤", stepId)))
		}
		//结束步骤，退出创建任务
		if pNextStep.Category == StepCat.END.Code {
			break
		}
		//创建下一个步骤的任务
		NewTaskByStep(pNextStep, pProcess, 1, form, tx)
		stepId = pNextStep.Id
	}

	//查询指定步骤之后的所有未启动的任务,并更新状态
	unstartTasks := Query(dto.TaskQueryDto{ProcessId: processId, StartStepId: startStepId, State: TaskState.UNSTART.Code}, tx)
	if len(unstartTasks) < 1 {
		FinishPassProcess(pProcess, tx)
	}
	//更新任务状态
	for _, unstartTask := range unstartTasks {
		//抄送任务自动完成
		if unstartTask.Category == StepCat.NOTIFY.Code {
			if unstartTask.State != TaskState.PASS.Code {
				unstartTask.State = TaskState.PASS.Code
				dao.SaveOrUpdate(&unstartTask, tx)
			}
			//只启动第一个审核任务
		} else if unstartTask.Category == StepCat.AUDIT.Code {
			if unstartTask.State == TaskState.UNSTART.Code {
				unstartTask.State = TaskState.STARTED.Code
				dao.SaveOrUpdate(&unstartTask, tx)
				break
			}
			//	没有抄送、审核任务，结束流程
		} else if unstartTask.Category == StepCat.END.Code { //结束步骤,结束流程
			FinishPassProcess(pProcess, tx)
			break
		}
	}
}

// 调用流程任务状态变更通知回调接口
func NotifyTasksStateChange(processId int, tx *gorm.DB) {
	url := config.Config.Notify.TaskStateChange
	if len(url) == 0 {
		log.Println("通知消息确认失败")
		log.Println("无效的流程任务状态变更通知回调地址")
		return
	}

	notifyVo := GetTaskStateChangeVo(processId, tx)
	m, err := CollectionUtil.Obj2map(notifyVo)
	if nil != err {
		log.Println("通知数据转map失败: %v", err)
		return
	}
	result := AjaxJson.AjaxJson{}
	RequestUtil.PostJson(url, m, &result)
	log.Println("接收通知服务端返回: %v", result)

	if CONSTANT.SUCESS_CODE != result.Code {
		log.Println("通知消息确认失败,返回:")
		log.Println(JsonUtil.Obj2PrettyJson(result))
	}
}

func FindPrevStep(processId int, targetStepId int, tx *gorm.DB) entity.Step {
	pProcess := dao.CheckById[entity.Process](processId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	prevStep := StepService.FindPrevStep(&pTemplate.RootStep, targetStepId)
	return *prevStep
}

func GetTaskStateChangeVo(processId int, tx *gorm.DB) vo.TaskStateChangeNotifyVo {
	notifyVo := vo.TaskStateChangeNotifyVo{}
	pProcess := dao.CheckById[entity.Process](processId, tx)
	notifyVo.Process = *pProcess

	tasks := TaskAssigneeDao.GetTasksByLastSubmitIndex(processId, tx)

	taskVos := []vo.TaskStateChangeNotifyTaskVo{}
	for _, v := range tasks {
		taskVo := vo.TaskStateChangeNotifyTaskVo{}
		taskVo.Task = v

		assignees := TaskAssigneeDao.GetLastSubmitAssigneesByTask(v.ProcessId, v.Id, tx)
		taskVo.Assignees = assignees

		taskVos = append(taskVos, taskVo)
	}
	notifyVo.Tasks = taskVos

	return notifyVo
}

// 查询流程实例的任务列表
func Query(taskQueryDto dto.TaskQueryDto, tx *gorm.DB) []entity.Task {
	return TaskDao.Query(taskQueryDto, tx)
}

func GetStartedTask(processId int, tx *gorm.DB) entity.Task {
	return TaskDao.GetStartedTask(processId, tx)
}
