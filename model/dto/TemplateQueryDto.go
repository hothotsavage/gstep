package dto

import PageDto "github.com/hothotsavage/gstep/util/db/page"

type TemplateQueryDto struct {
	PageDto.PageDto
	Id        int
	MouldId   int
	VersionId int
}
