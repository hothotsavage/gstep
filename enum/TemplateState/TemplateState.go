package TemplateState

import "github.com/hothotsavage/gstep/util/enum"

type TemplateState struct {
	enum.BaseEnum[string]
}

var RELEASE = TemplateState{}
var DRAFT = TemplateState{}

func init() {
	RELEASE.Code = "release"
	RELEASE.Title = "发布"

	DRAFT.Code = "draft"
	DRAFT.Title = "草稿"
}
