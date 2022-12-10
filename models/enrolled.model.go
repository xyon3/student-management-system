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

type StudentEnrolledCourses struct {
	CourseID    string `json:"courseID" gorm:"primaryKey;column:courseID"`
	Description string `json:"description"`
	Proctor     string `json:"proctor"`
	Day         string `json:"day"`
	StartTime   string `json:"startTime" gorm:"column:startTime"`
	EndTime     string `json:"endTime" gorm:"column:endTime"`
	RingDelay   string `json:"ringDelay" gorm:"column:ringDelay"`
}

func (Enrolled) TableName() string            { return "juncEnrolled" }
func (EnrolledRequestBody) TableName() string { return "juncEnrolled" }
