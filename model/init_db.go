package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sqlOJ/config"
)

var sysDB *gorm.DB
var exeDB *gorm.DB

func InitDB() {
	InitSysDB()
	InitExeDB()
}

func InitSysDB() {
	// configure parameters for mysql
	username := config.Username
	password := config.Password
	host := config.Host
	port := config.Port
	SysDBName := config.SysDBName
	//dsn := "root:spln13spln@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, SysDBName)
	var err error
	sysDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // print sql sentences
	})
	if err != nil { // fail to connect with mysql, need to use panic
		panic(err)
	}
	// create tables in database
	err = sysDB.AutoMigrate(&AdminAccount{}, &TeacherAccount{}, &StudentAccount{}, &ExerciseAssociation{},
		&ExerciseTable{}, &ExerciseContent{}, &SubmitHistory{}, &UserProblemStatus{}, &Contest{},
		&ContestClassAssociation{}, &ContestExerciseAssociation{}, &ScoreRecord{}, &ContestSubmission{}, &ContestExerciseStatus{})
	if err != nil {
		panic(err)
	}

	db, err := sysDB.DB()
	if err != nil {
		log.Fatalln("db connected error", err)
	}
	//db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
}

func InitExeDB() {
	username := config.Username
	password := config.Password
	host := config.Host
	port := config.Port
	ExeDBName := config.ExeDBName
	//dsn := "root:spln13spln@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, ExeDBName)
	var err error
	exeDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // print sql sentences
	})
	if err != nil { // fail to connect with mysql, need to use panic
		panic(err)
	}

	//err = execDB.AutoMigrate(&UserInfo{})
	//if err != nil {
	//	panic(err)
	//}

	db, err := exeDB.DB()
	if err != nil {
		log.Fatalln("db connected error", err)
	}
	//db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
}

func GetSysDB() *gorm.DB {
	return sysDB.Session(&gorm.Session{
		SkipDefaultTransaction: true, // 禁用默认事务
		PrepareStmt:            true, // 缓存预编译命令
	})
}

func GetExeDB() *gorm.DB {
	return exeDB.Session(&gorm.Session{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
}
