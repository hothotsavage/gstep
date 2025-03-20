package TaskService

import (
	"fmt"
	"github.com/gookit/goutil/strutil"
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
	"github.com/hothotsavage/gstep/service/TaskAssigneeService"
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
	pStep := StepService.FindStep(&pTemplate.RootStep, currentStepId, tx)

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
		pPrevBranchStep := StepService.FindPrevBranchStepWithNextStep(&pTemplate.RootStep, pStep.Id, tx)
		if nil != pPrevBranchStep.NextStep {
			return pPrevBranchStep.NextStep
		}
	}

	return nil
}

func NewTaskByStep(pStep *entity.Step, pProcess *entity.Process, form *map[string]any, tx *gorm.DB) *entity.Task {
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

	task := entity.Task{}
	task.ProcessId = pProcess.Id
	task.StepId = pStep.Id
	task.Title = pStep.Title
	task.Category = pStep.Category
	task.AuditMethod = pStep.AuditMethod
	task.Form = nil
	task.State = TaskState.UNSTART.Code
	task.Candidates = MakeCandidates(*pStep, form, tx)

	//重新创建启动任务，候选人为第一个启动任务的申请人
	if pStep.Category == StepCat.START.Code {
		firstTaskSubmitterUserId := TaskAssigneeDao.GetFirstTaskSubmitter(pProcess.Id, tx)
		if strutil.IsNotBlank(firstTaskSubmitterUserId) {
			task.Candidates = []string{firstTaskSubmitterUserId}
		}
	}

	dao.SaveOrUpdate(&task, tx)
	return &task
}

// 创建启动任务
func NewStartTask(pProcess *entity.Process, startUserId string, form *map[string]any, memo string, tx *gorm.DB) (*entity.Task, *entity.TaskAssignee) {
	//创建启动任务
	task := entity.Task{}
	task.ProcessId = pProcess.Id

	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	rootStep := pTemplate.RootStep

	//创建启动任务
	task.StepId = rootStep.Id
	task.Title = rootStep.Title
	task.Form = form
	task.Category = rootStep.Category
	task.State = TaskState.PASS.Code
	task.Candidates = MakeCandidates(rootStep, form, tx)
	task.AuditMethod = rootStep.AuditMethod
	if strutil.IsBlank(task.AuditMethod) {
		task.AuditMethod = AuditMethodCat.OR.Code
	}
	dao.SaveOrUpdate(&task, tx)

	//创建启动任务的执行人
	assignee := entity.TaskAssignee{}
	assignee.ProcessId = pProcess.Id
	assignee.StepId = task.StepId
	assignee.TaskId = task.Id
	assignee.UserId = startUserId
	assignee.State = TaskState.PASS.Code
	assignee.SubmitIndex = 1
	assignee.Form = form
	assignee.Memo = memo
	dao.SaveOrUpdate(&assignee, tx)

	return &task, &assignee
}

// 审核通过流程
func FinishPassProcess(pProcess *entity.Process, tx *gorm.DB) {
	pProcess.State = ProcessState.FINISH_PASS.Code
	finishTime := LocalTime.LocalTime(time.Now())
	pProcess.FinishedAt = &finishTime
	dao.SaveOrUpdate(pProcess, tx)
}

// 创建指定步骤之后的所有任务列表
func MakeTasks(processId int, startStepId int, form *map[string]any, tx *gorm.DB) []string {
	pProcess := dao.CheckById[entity.Process](processId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	pStartStep := StepService.FindStep(&pTemplate.RootStep, startStepId, tx)
	if pStartStep == nil {
		panic(ServerError.New(fmt.Sprintf("找不到流程步骤(stepId=%s)", startStepId)))
	}
	if pStartStep.Category == StepCat.END.Code {
		return []string{}
	}

	//创建第一个步骤的任务实例
	pStartTask := NewTaskByStep(pStartStep, pProcess, form, tx)
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
		//创建下一个步骤的任务实例
		NewTaskByStep(pNextStep, pProcess, form, tx)
		stepId = pNextStep.Id
	}

	//查询指定任务之后的所有未启动的任务,并更新状态
	unstartTasks := TaskDao.Query(dto.TaskQueryDto{ProcessId: processId, StartTaskId: pStartTask.Id, State: TaskState.UNSTART.Code}, tx)
	notifyUserIds := []string{}
	//更新任务状态
	for _, unstartTask := range unstartTasks {
		//抄送任务自动完成
		if unstartTask.Category == StepCat.NOTIFY.Code {
			if unstartTask.State != TaskState.PASS.Code {
				unstartTask.State = TaskState.PASS.Code
				dao.SaveOrUpdate(&unstartTask, tx)
				notifyUserIds = append(notifyUserIds, unstartTask.Candidates...)
				//创建抄送人列表
				MakeAssignees(unstartTask, form, tx)
			}
			//只启动第一个审核或开始任务
		} else if unstartTask.Category == StepCat.START.Code {
			if unstartTask.State == TaskState.UNSTART.Code {
				unstartTask.State = TaskState.STARTED.Code
				dao.SaveOrUpdate(&unstartTask, tx)
				notifyUserIds = append(notifyUserIds, unstartTask.Candidates...)
				break
			}
		} else if unstartTask.Category == StepCat.AUDIT.Code {
			if unstartTask.State == TaskState.UNSTART.Code {
				unstartTask.State = TaskState.STARTED.Code
				dao.SaveOrUpdate(&unstartTask, tx)
				notifyUserIds = append(notifyUserIds, unstartTask.Candidates...)
				break
			}
		}
	}

	return notifyUserIds
}

// 生成任务的通知消息文案
func MakeNotifyMessage(taskAssignee entity.TaskAssignee, tx *gorm.DB) string {
	name := ""
	if strutil.IsNotBlank(taskAssignee.UserId) {
		pUser := dao.CheckById[entity.User](taskAssignee.UserId, tx)
		if nil != pUser {
			name = pUser.Name
		}
	}

	action := ""
	pTask := dao.CheckById[entity.Task](taskAssignee.TaskId, tx)
	if nil != pTask {
		if pTask.Category == StepCat.START.Code {
			action = "申请了"
		} else if pTask.Category == StepCat.AUDIT.Code {
			action = "审核了"
		} else if pTask.Category == StepCat.NOTIFY.Code {
			action = "抄送了"
		}
	}

	bill := ""
	pProcess := dao.CheckById[entity.Process](taskAssignee.ProcessId, tx)
	if pProcess != nil {
		pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
		if pTemplate != nil {
			bill = pTemplate.Title
		}
	}

	return fmt.Sprintf("%s%s%s", name, action, bill)
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

// 检查指定步骤的候选人
func CheckCandidate(userId string, form *map[string]any, taskId int, tx *gorm.DB) {
	pTask := dao.CheckById[entity.Task](taskId, tx)
	//没有候选人名单，表示所有人都可提交，直接通过
	if len(pTask.Candidates) < 1 {
		return
	}

	for _, candidate := range pTask.Candidates {
		if userId == candidate {
			return
		}
	}

	//for _, v := range pTask.Candidates {
	//	if v.Category == CandidateCat.USER.Code {
	//		if userId == v.Value {
	//			return
	//		}
	//	} else if v.Category == CandidateCat.DEPARTMENT.Code {
	//		departments := DepartmentDao.GetGrandsonDepartments(v.Value)
	//		isIn := UserDao.IsUserInDepartments(userId, departments)
	//		if isIn {
	//			return
	//		}
	//	} else if v.Category == CandidateCat.FIELD.Code {
	//		formCandidate := (*form)[v.Value]
	//		if formCandidate == userId {
	//			return
	//		}
	//	}
	//}

	panic(ServerError.New("流程提交人不在候选人列表中"))
}

// 创建任务的候选人列表
func MakeCandidates(step entity.Step, form *map[string]any, tx *gorm.DB) []string {
	userIds := []string{}
	for _, c := range step.Candidates {
		if c.Category == CandidateCat.USER.Code {
			if strutil.IsNotBlank(c.Value) {
				userIds = append(userIds, c.Value)
			}
		} else if c.Category == CandidateCat.FIELD.Code {
			userId := (*form)[c.Value]
			if strutil.IsNotBlank(userId.(string)) {
				userIds = append(userIds, userId.(string))
			}
		} else if c.Category == CandidateCat.DEPARTMENT.Code {
			users := UserDao.GetGrandsonDepartmentUsers(c.Value, tx)
			for _, user := range users {
				if strutil.IsNotBlank(user.Id) {
					userIds = append(userIds, user.Id)
				}
			}
		}
	}
	return userIds
}

// 创建任务的执行人实例
func MakeAssignees(task entity.Task, form *map[string]any, tx *gorm.DB) {
	//记录提交人的排序序号
	submitIndex := 1
	for _, userId := range task.Candidates {
		assignee := entity.TaskAssignee{}
		assignee.ProcessId = task.ProcessId
		assignee.StepId = task.StepId
		assignee.TaskId = task.Id
		assignee.State = task.State
		assignee.SubmitIndex = submitIndex
		assignee.Form = form
		assignee.UserId = userId
		dao.SaveOrUpdate(&assignee, tx)
		submitIndex++
	}

	////记录提交人的排序序号
	//submitIndex := 1
	//for _, c := range task.Candidates {
	//	if c.Category == CandidateCat.USER.Code {
	//		assignee := entity.TaskAssignee{}
	//		assignee.ProcessId = task.ProcessId
	//		assignee.StepId = task.StepId
	//		assignee.TaskId = task.Id
	//		assignee.State = task.State
	//		assignee.SubmitIndex = submitIndex
	//		assignee.Form = form
	//		assignee.UserId = c.Value
	//		dao.SaveOrUpdate(&assignee)
	//	} else if c.Category == CandidateCat.FIELD.Code {
	//		assignee := entity.TaskAssignee{}
	//		assignee.ProcessId = task.ProcessId
	//		assignee.StepId = task.StepId
	//		assignee.TaskId = task.Id
	//		assignee.State = task.State
	//		assignee.SubmitIndex = submitIndex
	//		assignee.Form = form
	//		userId := (*form)[c.Value]
	//		assignee.UserId = userId.(string)
	//		dao.SaveOrUpdate(&assignee)
	//	} else if c.Category == CandidateCat.DEPARTMENT.Code {
	//		users := UserDao.GetGrandsonDepartmentUsers(c.Value)
	//		for _, UserDao := range users {
	//			assignee := entity.TaskAssignee{}
	//			assignee.ProcessId = task.ProcessId
	//			assignee.StepId = task.StepId
	//			assignee.TaskId = task.Id
	//			assignee.State = task.State
	//			assignee.SubmitIndex = submitIndex
	//			assignee.Form = form
	//			assignee.UserId = UserDao.Id
	//			dao.SaveOrUpdate(&assignee)
	//		}
	//	}
	//	submitIndex++
	//}
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

func GetTasksByProcessId(processId int, tx *gorm.DB) []entity.Task {
	tasks := []entity.Task{}
	tx.Raw("select * from task where process_id=? order by id asc", processId).Scan(&tasks)
	return tasks
}

func ToVO(task entity.Task, tx *gorm.DB) vo.TaskVO {
	aVO := vo.TaskVO{}
	aVO.Task = task
	assignees := TaskAssigneeDao.GetAssigneesByTaskId(task.Id, tx)
	assigneeVOS := []vo.TaskAssigneeVO{}
	for _, v := range assignees {
		assigneeVO := TaskAssigneeService.ToVO(v, tx)
		assigneeVOS = append(assigneeVOS, assigneeVO)
	}
	aVO.Assignees = assigneeVOS

	candidateUsers := []entity.User{}
	for _, v := range aVO.Candidates {
		user := dao.GetById[entity.User](v, tx)
		candidateUsers = append(candidateUsers, *user)
	}
	aVO.CandidateUsers = candidateUsers
	return aVO
}
