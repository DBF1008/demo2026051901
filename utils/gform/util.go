package gform

import (
	"fmt"
	"github.com/gohouse/t"
	"log"
	"math/rand"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"
)

func getRandomInt(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

func structForScan(u interface{}) []interface{} {
	val := reflect.Indirect(reflect.ValueOf(u))
	v := make([]interface{}, 0)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if val.Type().Field(i).Tag.Get(TAGNAME) != IGNORE {
			if valueField.CanAddr() {
				v = append(v, valueField.Addr().Interface())
			} else {
				//v[i] = valueField
				v = append(v, valueField)
			}
		}
	}
	return v
}

// StructToMap ...
func StructToMap(obj interface{}) map[string]interface{} {
	ty := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < ty.NumField(); i++ {
		data[ty.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func getTagName(structName interface{}, tagstr string) []string {
	tag := reflect.TypeOf(structName)
	if tag.Kind() == reflect.Ptr {
		tag = tag.Elem()
	}

	if tag.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := tag.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := tag.Field(i).Tag.Get(tagstr)
		if tagName != IGNORE {
			if tagName == "-" || tagName == "" {
				tagName = tag.Field(i).Name
			}
			result = append(result, tagName)
		}
	}
	return result
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func addQuotes(data interface{}, sep string) string {
	ret := t.New(data).String()
	ret = strings.Replace(ret, `\`, `\\`, -1)
	ret = strings.Replace(ret, `"`, `\"`, -1)
	ret = strings.Replace(ret, `'`, `\'`, -1)
	return fmt.Sprintf("%s%s%s", sep, ret, sep)
}

func inArray(needle, hystack interface{}) bool {
	nt := t.New(needle)
	for _, item := range t.New(hystack).Slice() {
		if strings.ToLower(nt.String()) == strings.ToLower(item.String()) {
			return true
		}
	}
	return false
}

func withLockContext(fn func()) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	fn()
}

func withRunTimeContext(closer func(), callback func(time.Duration)) {
	start := time.Now()
	closer()
	timeduration := time.Since(start)
	callback(timeduration)
}

func readFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path.Dir(filepath), os.ModePerm)
		file, _ = os.Create(filepath)
	}
	return file
}
