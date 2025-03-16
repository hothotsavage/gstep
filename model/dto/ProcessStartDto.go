package dto

type ProcessStartDto struct {
	MouldId int             `json:"mouldId"`
	Form    *map[string]any `json:"form"`
	UserId  string          `json:"userId"`
}
