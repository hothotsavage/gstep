package entity

type Position struct {
	Title string `json:"title" gorm:"primarykey"`
	Code  string `json:"code"`
}

func (e Position) TableName() string {
	return "UserDao"
}

func (e Position) GetId() any {
	return e.Code
}
