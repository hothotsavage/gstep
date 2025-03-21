package StepService

import (
	"fmt"
	"github.com/gookit/goutil/strutil"
	"github.com/hothotsavage/gstep/dao/DepartmentDao"
	"github.com/hothotsavage/gstep/dao/UserDao"
	"github.com/hothotsavage/gstep/enum/CandidateCat"
	"github.com/hothotsavage/gstep/enum/StepCat"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func GetStepByTemplateId(templateId int, stepId int, tx *gorm.DB) *entity.Step {
	template := dao.CheckById[entity.Template](templateId, tx)
	pStep := FindStep(&template.RootStep, stepId, tx)
	return pStep
}

func FindStep(pParentStep *entity.Step, stepId int, tx *gorm.DB) *entity.Step {
	if nil == pParentStep {
		return nil
	}

	if pParentStep.Id == stepId {
		return pParentStep
	}

	if nil != pParentStep.NextStep && pParentStep.NextStep.Id == stepId {
		return pParentStep.NextStep
	}

	for _, v := range pParentStep.BranchSteps {
		if v.Id == stepId {
			return v
		}
	}

	if nil != pParentStep {
		aStep := FindStep(pParentStep.NextStep, stepId, tx)
		if nil != aStep {
			return aStep
		}
	}

	for _, v := range pParentStep.BranchSteps {
		pFindOne := FindStep(v, stepId, tx)
		if nil != pFindOne {
			return pFindOne
		}
	}

	return nil
}

// 递归查找指定步骤的前一个步骤
// pParentStep 父步骤
// targetStepId 查询的步骤id
func FindPrevStep(pParentStep *entity.Step, targetStepId int, tx *gorm.DB) *entity.Step {
	if nil == pParentStep {
		return nil
	}

	if nil != pParentStep.NextStep && pParentStep.NextStep.Id == targetStepId {
		return pParentStep
	}

	for _, v := range pParentStep.BranchSteps {
		if v.Id == targetStepId {
			return pParentStep
		}
	}

	if nil != pParentStep.NextStep {
		pNextStep := FindPrevStep(pParentStep.NextStep, targetStepId, tx)
		if nil != pNextStep {
			return pNextStep
		}
	}

	for _, v := range pParentStep.BranchSteps {
		pFindOne := FindPrevStep(v, targetStepId, tx)
		if nil != pFindOne {
			return pFindOne
		}
	}

	return nil
}

// 递归查找前一个分支步骤
func FindPrevBranchStepWithNextStep(pRootStep *entity.Step, targetStepId int, tx *gorm.DB) *entity.Step {
	pPrevStep := FindPrevStep(pRootStep, targetStepId, tx)

	if nil == pPrevStep {
		return nil
	}

	if pPrevStep.Category == StepCat.BRANCH.Code && nil != pPrevStep.NextStep && pPrevStep.NextStep.Id != 0 {
		return pPrevStep
	}

	pPrevPrevStep := FindPrevBranchStepWithNextStep(pRootStep, pPrevStep.Id, tx)
	return pPrevPrevStep
}

// 查找前一个审核步骤
func FindPrevAuditStep(pRootStep *entity.Step, beginStepId int, tx *gorm.DB) *entity.Step {
	fromStepId := beginStepId
	for {
		pPrevStep := FindPrevStep(pRootStep, fromStepId, tx)
		if nil == pPrevStep || pPrevStep.Id == 0 {
			return nil
		}
		if StepCat.IsContainAudit(pPrevStep.Category) {
			return pPrevStep
		}

		fromStepId = pPrevStep.Id
	}
}

// 查询前面所有审核步骤列表
func FindPrevAuditSteps(pRootStep *entity.Step, beginStepId int, tx *gorm.DB) []entity.Step {
	auditSteps := []entity.Step{}
	fromStepId := beginStepId
	for {
		pPrevStep := FindPrevAuditStep(pRootStep, fromStepId, tx)
		if nil == pPrevStep {
			return auditSteps
		}
		if StepCat.IsContainAudit(pPrevStep.Category) {
			auditSteps = append(auditSteps, *pPrevStep)
		}

		fromStepId = pPrevStep.Id
	}
}

// 查询之前到指定步骤的审核列表
func FindPrevAuditStepsByEndId(pRootStep *entity.Step, beginStepId int, endStepId int, tx *gorm.DB) []entity.Step {
	auditpSteps := []entity.Step{}
	fromStepId := beginStepId
	for {
		pPrevStep := FindPrevAuditStep(pRootStep, fromStepId, tx)

		if nil == pPrevStep {
			return auditpSteps
		}

		if StepCat.IsContainAudit(pPrevStep.Category) {
			auditpSteps = append(auditpSteps, *pPrevStep)
		}

		if endStepId == pPrevStep.Id {
			return auditpSteps
		}

		fromStepId = pPrevStep.Id
	}
}

// 检查指定步骤的候选人
func CheckStepCandidate(userId string, pTaskForm *map[string]any, templateId int, stepId int, tx *gorm.DB) {
	pTemplate := dao.CheckById[entity.Template](templateId, tx)
	pStep := FindStep(&pTemplate.RootStep, stepId, tx)
	//没有候选人名单，表示所有人都可提交，直接通过
	if len(pStep.Candidates) < 1 {
		return
	}

	candidateStr := ""
	for _, v := range pStep.Candidates {
		if v.Category == CandidateCat.USER.Code {
			candidateStr += v.Value + ","

			if userId == v.Value {
				return
			}
		} else if v.Category == CandidateCat.DEPARTMENT.Code {
			candidateStr += v.Value + ","

			departments := DepartmentDao.GetGrandsonDepartments(v.Value, tx)
			isIn := UserDao.IsUserInDepartments(userId, departments, tx)
			if isIn {
				return
			}
		} else if v.Category == CandidateCat.FIELD.Code {
			candidateStr += v.Value + ","

			formValue := (*pTaskForm)[v.Value]
			if nil == formValue {
				panic(ServerError.New(fmt.Sprintf("任务表单中没有字段(%s)", v.Value)))
			}
			switch formValue.(type) {
			case string:
				formCandidateUserId := formValue.(string)
				if strutil.IsBlank(userId) {
					panic(ServerError.New(fmt.Sprintf("任务表单中字段(%s)值为空", v.Value)))
				}
				if formCandidateUserId == userId {
					return
				}
			case []interface{}:
				formCandidateUserIds := cast.ToStringSlice(formValue)
				if nil == formCandidateUserIds || len(formCandidateUserIds) < 1 {
					panic(ServerError.New(fmt.Sprintf("任务表单中字段(%s)值为空", v.Value)))
				}
				for _, formUserId := range formCandidateUserIds {
					if formUserId == userId {
						return
					}
				}
			default:
				panic(ServerError.New(fmt.Sprintf("任务表单中字段(%s)类型错误", v.Value)))
			}
		}
	}

	//删除最后一个逗号
	if len(candidateStr) > 0 {
		candidateStr = candidateStr[:len(candidateStr)-1]
	}
	panic(ServerError.New(fmt.Sprintf("流程提交人(userId=%s)不在步骤(%s)候选人(%s)中", userId, pStep.Title, candidateStr)))
}

// 候选人条数
func CandidateCount(taskId int, tx *gorm.DB) int {
	pTask := dao.CheckById[entity.Task](taskId, tx)
	pProcess := dao.CheckById[entity.Process](pTask.ProcessId, tx)
	pTemplate := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	pStep := FindStep(&pTemplate.RootStep, pTask.StepId, tx)
	return len(pStep.Candidates)
}
