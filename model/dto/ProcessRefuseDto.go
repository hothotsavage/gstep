package dto

type ProcessRefuseDto struct {
	ProcessId  int
	Form       *map[string]any
	UserId     string
	PrevStepId int
}
