package utils

import "sqlOJ/model"

func QueryUsername(userID int64, userType int64) string {
	var username string
	if userType == 1 { // 学生
		username = model.NewStudentAccountFlow().QueryStudentUsernameByUserID(userID)
	} else if userType == 2 { // 教师
		username = model.NewTeacherAccountFlow().QueryTeacherUsernameByUserID(userID)
	} else { // 管理员
		username = model.NewAdminAccountFlow().QueryAdminUsernameByUserID(userID)
	}
	return username
}
