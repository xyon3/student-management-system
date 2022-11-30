package models

type Student struct {
	StudID     string `json:"studID"  gorm:"primaryKey;column:studID"`
	Name       string `json:"name"`
	ProfileImg string `json:"profileImg" gorm:"column:profileImg"`
	Hash       string `json:"hash"`
	RingDelay  string `json:"ringDelay" gorm:"column:ringDelay"`
}

type StudentLoginBody struct {
	StudID string `json:"studID"  gorm:"column:studID"`
	Hash   string `json:"hash"`
}

func (Student) TableName() string          { return "tblStudent" }
func (StudentLoginBody) TableName() string { return "tblStudent" }
