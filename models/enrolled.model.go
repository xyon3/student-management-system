package models

type Enrolled struct {
	StudID   string `json:"studID" gorm:"column:studID"`
	CourseID string `json:"courseID"  gorm:"column:courseID"`
	ID       int64  `json:"id"`
}

type EnrolledRequestBody struct {
	StudID   string `json:"studID" gorm:"column:studID"`
	CourseID string `json:"courseID" gorm:"column:courseID"`
}

func (Enrolled) TableName() string            { return "juncEnrolled" }
func (EnrolledRequestBody) TableName() string { return "juncEnrolled" }
