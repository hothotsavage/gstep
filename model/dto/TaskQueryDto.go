package dto

type TaskQueryDto struct {
	ProcessId   int    `json:"processId"`
	StartStepId int    `json:"startStepId"`
	State       string `json:"state"`
	Category    string `json:"category"`
}
