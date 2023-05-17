package cache

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// ContestStudentCache 将参与该竞赛的学生ID缓存在Redis中
// 设置生效时间(单独key-value)为竞赛开始时间, ExpireAt endAt
func ContestStudentCache(contestID int64, studentIDList []int64, beginAt, endAt time.Time) error {
	contestIDStr := strconv.FormatInt(contestID, 10)
	setName := "contest_forbid_student:" + contestIDStr
	for _, studentID := range studentIDList {
		if err := rdb.SAdd(ctx, setName, studentID).Err(); err != nil {
			log.Println(err)
			return errors.New("缓存参与竞赛学生ID错误")
		}
	}
	if err := rdb.ExpireAt(ctx, setName, endAt).Err(); err != nil {
		log.Println(err)
		return errors.New("设置参与名单Set过期时间错误")
	}
	// 此键值不仅可以标记竞赛开始时间，还可以跟此键值是否存在来判断该竞赛是否结束
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
// !!!!!!!该函数已弃用!!!!!!!!
func ExerciseContestCache(contestID, exerciseID int64) error {
	exerciseIDStr := strconv.FormatInt(exerciseID, 10)
	setName := "exercise_contest:" + exerciseIDStr
	if err := rdb.SAdd(ctx, setName, contestID).Err(); err != nil {
		log.Println(err)
		return errors.New("缓存引用竞赛的题目时错误")
	}
	return nil
}

// GetExerciseSetMember 获取exerciseID对应的Set中所有contestID
func GetExerciseSetMember(exerciseIDStr string) ([]string, error) {
	setName := "exercise_contest:" + exerciseIDStr
	members, err := rdb.SMembers(ctx, setName).Result()
	if err != nil {
		log.Println(err)
		return nil, errors.New("查找引用题目的竞赛时错误")
	}
	return members, nil
}

// CheckRequestValidInContest 检查contestIDStr对应的Set是否已经过期
// 首先判断contestIDStr对应的Set是否过期，再判断userID是否在contestID对应的Set中
// 返回值规定: 0->错误; 1->该键值不存在; 2->集合中存在; 3->集合中不存在或未到时间
func CheckRequestValidInContest(userID int64, contestIDStr string) (int, error) {
	keyName := "contest_valid_time_name:" + contestIDStr
	exists, err := rdb.Exists(ctx, keyName).Result() // 首先查询contest_valid_time是否存在
	if err != nil {
		log.Println(err)
		return 0, errors.New("查询竞赛是否存在时错误")
	}
	if exists != 1 { // 键值不存在，表示竞赛已结束
		return 1, nil
	}
	contestValidTimeStr, err := rdb.Get(ctx, keyName).Result() // 获取缓存中竞赛开始时间
	if err != nil {
		log.Println(err)
		return 0, errors.New("查询竞赛是否存在时错误")
	}
	contestValidTime, _ := strconv.ParseInt(contestValidTimeStr, 10, 64) // 获取到竞赛开始时间的Unix值
	currentTime := time.Now().Unix()                                     // 获取当前时间Unix
	if currentTime < contestValidTime {                                  // 若当前时间未到竞赛开始时间
		return 3, nil
	}
	setName := "contest_forbid_student:" + contestIDStr
	//userIDStr := strconv.FormatInt(userID, 10)
	fmt.Println(setName)
	exist, err := rdb.SIsMember(ctx, setName, userID).Result() // 判断当前用户ID是否在Set中
	fmt.Println(exist)
	if err != nil {
		log.Println(err)
		return 0, errors.New("查询用户是否被禁止访问时错误")
	}
	if !exist { // 如果集合中不存在
		return 3, nil
	}
	// 集合中存在
	return 2, nil
}

// DeleteContestIDInExercise 在exerciseIDStr对应的Set中删除contestIDStr
func DeleteContestIDInExercise(exerciseIDStr, contestIDStr string) error {
	setName := "exercise_contest:" + exerciseIDStr
	if err := rdb.SRem(ctx, setName, contestIDStr).Err(); err != nil {
		log.Println(err)
		return errors.New("删除引用题目的竞赛集合错误")
	}
	return nil
}
