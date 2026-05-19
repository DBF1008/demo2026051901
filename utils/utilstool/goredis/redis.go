package goredis

import (
	"context"
	"fmt"
	"time"
)


var client = GetRedisClient()
var ctx = context.Background()


func Set(key, value string) bool {
	result, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result == "OK"
}

func SetEX(key, value string, ex time.Duration) bool {
	result, err := client.Set(ctx, key, value, ex).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result == "OK"
}

func Get(key string) (bool, string) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, result
}

func GetSet(key, value string) (bool, string) {
	oldValue, err := client.GetSet(ctx, key, value).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, oldValue
}

func Incr(key string) int64 {
	val, err := client.Incr(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func IncrBy(key string, incr int64) int64 {
	val, err := client.IncrBy(ctx, key, incr).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func IncrByFloat(key string, incrFloat float64) float64 {
	val, err := client.IncrByFloat(ctx, key, incrFloat).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func Decr(key string) int64 {
	val, err := client.Decr(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func DecrBy(key string, incr int64) int64 {
	val, err := client.DecrBy(ctx, key, incr).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func Del(key string) bool {
	result, err := client.Del(ctx, key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

func Expire(key string, ex time.Duration) bool {
	result, err := client.Expire(ctx, key, ex).Result()
	if err != nil {
		return false
	}
	return result
}


func LPush(key string, date ...interface{}) int64 {
	result, err := client.LPush(ctx, key, date).Result()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func RPush(key string, date ...interface{}) int64 {
	result, err := client.RPush(ctx, key, date).Result()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func LPop(key string) (bool, string) {
	val, err := client.LPop(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, val
}

func RPop(key string) (bool, string) {
	val, err := client.RPop(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, val
}

func LIndex(key string, index int64) (bool, string) {
	val, err := client.LIndex(ctx, key, index).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, val
}

func LLen(key string) int64 {
	val, err := client.LLen(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func LRange(key string, start, stop int64) []string {
	vales, err := client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		fmt.Println(err)
	}
	return vales
}

func LRem(key string, count int64, data interface{}) bool {
	_, err := client.LRem(ctx, key, count, data).Result()
	if err != nil {
		fmt.Println(err)
	}
	return true
}

func LInsert(key string, pivot int64, data interface{}) bool {
	err := client.LInsert(ctx, key, "after", pivot, data).Err()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}


func SAdd(key string, data ...interface{}) bool {
	err := client.SAdd(ctx, key, data).Err()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func SCard(key string) int64 {
	size, err := client.SCard(ctx, "key").Result()
	if err != nil {
		fmt.Println(err)
	}
	return size
}

func SIsMember(key string, data interface{}) bool {
	ok, err := client.SIsMember(ctx, key, data).Result()
	if err != nil {
		fmt.Println(err)
	}
	return ok
}

func SMembers(key string) []string {
	es, err := client.SMembers(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return es
}

func SRem(key string, data ...interface{}) bool {
	_, err := client.SRem(ctx, key, data).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func SPopN(key string, count int64) []string {
	vales, err := client.SPopN(ctx, key, count).Result()
	if err != nil {
		fmt.Println(err)
	}
	return vales
}


func HSet(key, field, value string) bool {
	err := client.HSet(ctx, key, field, value).Err()
	if err != nil {
		return false
	}
	return true
}

func HGet(key, field string) string {
	val, err := client.HGet(ctx, key, field).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func HMGet(key string, fields ...string) []interface{} {
	vales, err := client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		panic(err)
	}
	return vales
}

func HGetAll(key string) map[string]string {
	data, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func HKeys(key string) []string {
	fields, err := client.HKeys(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return fields
}

func HLen(key string) int64 {
	size, err := client.HLen(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return size
}

func HMSet(key string, data map[string]interface{}) bool {
	result, err := client.HMSet(ctx, key, data).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result
}

func HSetNX(key, field string, value interface{}) bool {
	result, err := client.HSetNX(ctx, key, field, value).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result
}

func HDel(key string, fields ...string) bool {
	_, err := client.HDel(ctx, key, fields...).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func HExists(key, field string) bool {
	result, err := client.HExists(ctx, key, field).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result
}
