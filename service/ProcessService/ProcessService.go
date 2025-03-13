package ProcessService

import (
	"github.com/hothotsavage/gstep/dao/TemplateDao"
	"github.com/hothotsavage/gstep/enum/ProcessState"
	"github.com/hothotsavage/gstep/enum/StepCat"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/TaskService"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func Start(dto *dto.ProcessStartDto, tx *gorm.DB) int {
	process := entity.Process{}
	copier.Copy(process, dto)

	//创建流程
	pTemplate := TemplateDao.GetLatestVersionByTemplateId(dto.TemplateId, tx)
	if nil == pTemplate {
		panic(ServerError.New("无效的模板"))
	}
	process.TemplateId = pTemplate.Id
	process.StartUserId = dto.StartUserId
	process.State = ProcessState.STARTED.Code
	dao.SaveOrUpdate(&process, tx)

	//创建启动任务
	pTask := TaskService.NewStartTask(&process, dto.StartUserId, dto.Form, tx)

	//启动下一步
	for {
		pNextStep := TaskService.GetNextStep(pTask.StepId, pTemplate, pTask.Form, tx)
		if nil == pNextStep || 0 == pNextStep.Id {
			panic(ServerError.New("找不到流程启动步骤的下一个步骤"))
		}
		//创建新任务
		if pNextStep.Category != StepCat.END.Code {
			pTask = TaskService.NewTaskByStep(pNextStep, &process, 1, pTask.Form, tx)
		}

		//审核任务,退出
		if pNextStep.Category == StepCat.AUDIT.Code {
			break
		} else if pNextStep.Category == StepCat.END.Code { //结束步骤,结束流程
			TaskService.FinishPassProcess(&process, tx)
			break
		}
	}

	return process.Id
}
