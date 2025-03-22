package MouldService

import (
	"fmt"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/db/dao"
	"gorm.io/gorm"
)

func Save(pDetail *entity.Mould, tx *gorm.DB) {
	// 根据条件批量更新
	result := tx.Model(&entity.Template{}).Where("mould_id = ?", pDetail.Id).Updates(entity.Template{Title: pDetail.Title})
	if result.Error != nil {
		fmt.Println("更新流程模板(mould_id=%d)失败 ", pDetail.Id, result.Error)
	}

	dao.SaveOrUpdate(pDetail, tx)
}
