package dto

type ProcessStartDto struct {
	TemplateId int             `json:"templateId"`
	Form       *map[string]any `json:"form"`
	UserId     string          `json:"userId"`
}
