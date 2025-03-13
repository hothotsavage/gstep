package ProcessHandler

import (
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/ProcessService"
	"github.com/hothotsavage/gstep/service/TaskService"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func Start(writer http.ResponseWriter, request *http.Request) {
	requestDto := dto.ProcessStartDto{}
	RequestParsUtil.Body2dto(request, &requestDto)

	tx := DbUtil.GetTx()
	dao.CheckById[entity.User](requestDto.StartUserId, tx)
	//创建流程及启动任务
	id := ProcessService.Start(&requestDto, tx)

	//任务状态变更通知
	TaskService.NotifyTasksStateChange(id, tx)

	tx.Commit()

	AjaxJson.SuccessByData(id).Response(writer)
}
