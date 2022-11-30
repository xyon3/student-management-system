package models

type Course struct {
	CourseID    string `json:"courseID" gorm:"primaryKey;column:courseID"`
	Description string `json:"description"`
	Day         string `json:"day"`
	Time        string `json:"_time" gorm:"column:_time"`
}

func (Course) TableName() string { return "tblCourse" }
