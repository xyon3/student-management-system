package models

type Course struct {
	CourseID    string `json:"courseID" gorm:"primaryKey;column:courseID"`
	Description string `json:"description"`
	Proctor     string `json:"proctor"`
	Day         string `json:"day"`
	Room        string `json:"roomLoc" gorm:"column:roomLoc"`
	StartTime   string `json:"startTime" gorm:"column:startTime"`
	EndTime     string `json:"endTime" gorm:"column:endTime"`
}

func (Course) TableName() string { return "tblCourse" }
