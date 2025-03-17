package dto

type ProcessRefuseDto struct {
	ProcessId  int             `json:"processId"`
	Form       *map[string]any `json:"form"`
	Memo       string          `json:"memo"`
	UserId     string          `json:"userId"`
	PrevStepId int             `json:"prevStepId"`
}
