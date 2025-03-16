package TaskHandler

import (
	"github.com/hothotsavage/gstep/dao/TaskDao"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func Pending(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TaskPendingDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tasks, total := TaskDao.QueryMyPendingTasks(dto.UserId, DbUtil.Db)
	AjaxJson.SuccessByPagination(*tasks, total).Response(writer)
}

func Query(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TaskQueryDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tasks := TaskDao.Query(dto, DbUtil.Db)
	AjaxJson.SuccessByData(tasks).Response(writer)
}
