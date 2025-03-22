package entity

type Step struct {
	Id          int                  `json:"id"`
	Title       string               `json:"title"`
	Category    string               `json:"category"`
	Candidates  []Candidate          `json:"candidates"`
	Auth        map[string]FieldAuth `json:"auth" gorm:"serializer:json"`
	Expression  string               `json:"expression"`
	AuditMethod string               `json:"auditMethod"`
	BranchSteps []*Step              `json:"branchSteps"`
	NextStep    *Step                `json:"nextStep"`
}
