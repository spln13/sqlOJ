package cache

import (
	"errors"
	"log"
	"time"
)

// CheckEmailCodeSendTimeValid 验证请求时间是否大于cache中下次可发送时间
func CheckEmailCodeSendTimeValid(email string) (bool, error) {
	nextSendTimeKey := email + "_next_send_time"
	exists, err := rdb.Exists(ctx, email).Result()
	if err != nil {
		log.Println(err)
		return false, errors.New("查询缓存key错误")
	}
	if exists == 0 { // 不存在该键
		return true, nil // 合法
	}
	nextSendTime, err := rdb.Get(ctx, nextSendTimeKey).Int64()
	if err != nil {
		log.Println(err)
		return false, errors.New("检查发送时间错误")
	}
	if nextSendTime > time.Now().Unix() { // 重新请求少于1分钟
		return false, nil
	}
	return true, nil
}

// EmailCodeCache 将email以及生成的验证码保存到cache中，并设置下次可以发送的时间为一分钟后。
func EmailCodeCache(email string, code int, expire int) error {
	err := rdb.Set(ctx, email, code, time.Duration(expire)*time.Minute).Err()
	if err != nil {
		log.Println(err)
		return errors.New("缓存验证码失败")
	}
	nextSendTime := time.Now().Add(time.Minute).Unix() // 设置下次可以发送验证码的时间
	nextSendTimeKey := email + "_next_send_time"
	err = rdb.Set(ctx, nextSendTimeKey, nextSendTime, time.Duration(expire)*time.Minute).Err()
	if err != nil {
		log.Println(err)
		return errors.New("缓存验证码下次可发送时间错误")
	}
	return nil
}

func VerifyEmailCode(email string, code string) (bool, error) {
	cachedCode, err := rdb.Get(ctx, email).Result()
	if err != nil {
		log.Println(err)
		return false, errors.New("获取缓存验证码错误")
	}

	if err != nil {
		return false, errors.New("转换验证码格式错误")
	}

	if cachedCode != code {
		return false, errors.New("验证码错误")
	}

	err = rdb.Del(ctx, email).Err()
	if err != nil {
		log.Println(err)
		return false, errors.New("清除验证码缓存错误")
	}
	nextSendTimeKey := email + "_next_send_time"
	err = rdb.Del(ctx, nextSendTimeKey).Err()
	if err != nil {
		log.Println(err)
		return false, errors.New("清除验证码过期时间缓存错误")
	}

	return true, nil
}
