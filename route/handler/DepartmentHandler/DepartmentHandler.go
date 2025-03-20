package DepartmentHandler

import (
	"github.com/hothotsavage/gstep/ctx"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/service/DepartmentService"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func GetChildDepartments(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DepartmentQueryChildDto{}
	RequestParsUtil.Body2dto(request, &dto)
	tx := ctx.GetTx(request)

	childDepartments := DepartmentService.GetChildDepartments(dto, tx)

	AjaxJson.SuccessByData(childDepartments).Response(writer)
}

func GetUsers(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DepartmentQueryUsersDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := ctx.GetTx(request)

	users := DepartmentService.GetDepartmentUsers(dto, tx)

	AjaxJson.SuccessByData(users).Response(writer)
}
