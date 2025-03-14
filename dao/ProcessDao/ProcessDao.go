package ProcessDao

import (
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/model/vo"
	"github.com/hothotsavage/gstep/util/db/dao"
	"gorm.io/gorm"
)

func ToVo(pProcess *entity.Process, tx *gorm.DB) vo.ProcessVo {
	aVo := vo.ProcessVo{}
	if nil == pProcess {
		return aVo
	}

	aVo.Process = *pProcess
	template := dao.CheckById[entity.Template](pProcess.TemplateId, tx)
	aVo.Template = *template

	return aVo
}
