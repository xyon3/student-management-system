package models

type Enrolled struct {
	StudID   string `json:"studID"`
	CourseID string `json:"courseID"`
	ID       int64  `json:"id"`
}

func (Enrolled) TableName() string { return "juncEnrolled" }
