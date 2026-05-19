package common

import (
	"encoding/json"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Table struct {
}

func init() {
	fpath := Table{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Table) Weigh(context *gin.Context) {
	body, _ := io.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	bids, _ := json.Marshal(&ids)
	var ids_arr []interface{}
	var ids_arr_int []int64
	json.Unmarshal([]byte(bids), &ids_arr)
	json.Unmarshal([]byte(bids), &ids_arr_int)
	// ids_arr := strings.Split(ids, `,`)
	_changeid := parameter["changeid"].(float64)
	changeid := int64(_changeid)
	field := parameter["field"].(string)
	tablename := parameter["table"].(string)
	_pid := parameter["pid"].(float64)
	pid := int64(_pid)
	orderway := parameter["orderway"].(string)
	prikey := parameter["prikey"].(string)
	if _, ok := parameter["pid"]; ok {
		// var hasids []map[string]interface{}
		list_id, _ := model.DB().Table(tablename).WhereIn(prikey, ids_arr).Where("pid", pid).Pluck("id")
		list_int := list_id.([]interface{})
		list_intb := make([]int64, len(list_int))
		for i := range list_int {
			list_intb[i] = list_int[i].(int64)
		}
		ids_arr_int = intersect(list_intb, ids_arr_int)
	}
	winids_base, _ := json.Marshal(&ids_arr_int)
	var winids []interface{}
	_ = json.Unmarshal(winids_base, &winids)
	list, _ := model.DB().Table(tablename).WhereIn(prikey, winids).Fields(prikey + "," + field).Order(field + " " + orderway).Get()
	var sour []int64
	weighdata := make(map[int64]int64)
	for _, v := range list {
		sour = append(sour, v[prikey].(int64))
		weighdata[v[prikey].(int64)] = v[field].(int64)
	}
	position := array_search(changeid, ids_arr_int)
	desc_id := sour[position]
	change_id := changeid
	temp := difference(ids_arr_int, sour)
	for k, v := range temp {
		var offset int64
		if v == change_id {
			offset = desc_id
		} else {
			if change_id == temp[0] {
				nk := k + 1
				if len(temp) > nk {
					offset = temp[nk]
				} else {
					offset = change_id
				}
			} else {
				nk := k - 1
				if nk >= 0 {
					offset = temp[nk]
				} else {
					offset = change_id
				}
			}
		}
		model.DB().Table(tablename).Where(prikey, v).Data(map[string]interface{}{field: weighdata[offset]}).Update()
	}
	results.Success(context, "排序成功！", sour, desc_id)
}

func array_search(changeid int64, arr []int64) int {
	for k, v := range arr {
		if v == changeid {
			return k
		}
	}
	return -1
}

func intersect(nums1 []int64, nums2 []int64) []int64 {
	m := make(map[int64]int64)
	var arr []int64
	for _, v := range nums1 {
		m[v]++
	}
	for _, v := range nums2 {
		times, ok := m[v]
		if ok && times > 0 {
			arr = append(arr, v)
			m[v]--
		}
	}
	return arr
}

func difference(slice1, slice2 []int64) []int64 {
	var arr []int64
	for k, v := range slice1 {
		for key, value := range slice2 {
			if k == key && v != value {
				arr = append(arr, v)
			}
		}
	}
	if len(slice1) > len(slice2) {
		sn := len(slice2)
		n_arr := slice1[sn:]
		arr = ArrayMerge(arr, n_arr)
	}
	return arr
}

func ArrayMerge(ss ...[]int64) []int64 {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]int64, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}
