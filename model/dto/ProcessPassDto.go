package dto

type ProcessPassDto struct {
	ProcessId int             `json:"processId"`
	Form      *map[string]any `json:"form"`
	UserId    string          `json:"userId"`
}
