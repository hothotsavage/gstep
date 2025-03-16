package dto

type TaskQueryDto struct {
	ProcessId   int    `json:"processId"`
	StartTaskId int    `json:"startTaskId"`
	State       string `json:"state"`
	Category    string `json:"category"`
}
