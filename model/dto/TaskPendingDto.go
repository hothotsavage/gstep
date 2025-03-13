package dto

import PageDto "github.com/hothotsavage/gstep/util/db/page"

type TaskPendingDto struct {
	PageDto.PageDto
	UserId string
}
