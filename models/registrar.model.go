package models

type Registrar struct {
	RegID string `json:"regID" gorm:"primaryKey;column:regID"`
	Name  string `json:"name"`
	Hash  string `json:"hash"`
}

type RegistrarLoginBody struct {
	RegID string `json:"regID" gorm:"column:regID"`
	Hash  string `json:"hash"`
}

func (Registrar) TableName() string          { return "tblRegistrar" }
func (RegistrarLoginBody) TableName() string { return "tblRegistrar" }
