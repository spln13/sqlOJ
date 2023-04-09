package cache

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"
)

// ContestForbidStudentCache 将参与该竞赛的学生ID缓存在Redis中
// 设置生效时间(单独key-value)为竞赛开始时间, ExpireAt endAt
func ContestForbidStudentCache(contestID int64, studentIDList []int64, beginAt, endAt time.Time) error {
	byteData, err := json.Marshal(studentIDList)
	if err != nil {
		log.Println(err)
		return errors.New("缓存时转换学生ID类型错误")
	}
	contestIDStr := strconv.FormatInt(contestID, 10)
	setName := "contest_forbid_student:" + contestIDStr
	if err := rdb.SAdd(ctx, setName, byteData).Err(); err != nil {
		log.Println(err)
		return errors.New("缓存参与竞赛学生ID错误")
	}
	if err := rdb.ExpireAt(ctx, setName, endAt).Err(); err != nil {
		log.Println(err)
		return errors.New("设置禁止名单Set过期时间错误")
	}
	contestValidTimeName := "contest_valid_time_name:" + contestIDStr
	beginAtTimeStamp := beginAt.Unix()
	if err := rdb.Set(ctx, contestValidTimeName, beginAtTimeStamp, 0).Err(); err != nil {
		log.Println(err)
		return errors.New("设置竞赛生效时间错误")
	}
	if err := rdb.ExpireAt(ctx, contestValidTimeName, endAt).Err(); err != nil {
		log.Println(err)
		return errors.New("设置竞赛生效时间的过期时间错误")
	}
	return nil
}

// ExerciseContestCache 向exerciseID的Set中插入contestID
// 该Set不设置过期时间
func ExerciseContestCache(contestID, exerciseID int64) error {
	exerciseIDStr := strconv.FormatInt(exerciseID, 10)
	setName := "exercise_contest:" + exerciseIDStr
	if err := rdb.SAdd(ctx, setName, contestID).Err(); err != nil {
		log.Println(err)
		return errors.New("缓存引用竞赛的题目时错误")
	}
	return nil
}
