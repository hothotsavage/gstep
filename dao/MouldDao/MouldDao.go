package MouldDao

import (
	"fmt"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	PageDto "github.com/hothotsavage/gstep/util/db/page"
	pagination "github.com/hothotsavage/gstep/util/net/page"
	"gorm.io/gorm"
)

func Query(dto PageDto.PageDto, tx *gorm.DB) *pagination.Pagination[entity.Mould] {
	var entities []entity.Mould
	offset := (dto.Page - 1) * dto.Limit
	err := tx.Limit(dto.Limit).Offset(offset).Order("created_at desc").Find(&entities).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询流程模板失败 %s", err)))
	}

	var total int64
	tx.Model(&entity.Mould{}).Count(&total)

	res := pagination.Pagination[entity.Mould]{}
	res.List = entities
	res.Total = int(total)
	return &res
}
