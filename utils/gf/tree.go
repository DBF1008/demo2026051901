package gf

import "gofly/utils/gform"

func GetTreeArray(num []gform.Data, pid int64, itemprefix string) []gform.Data {
	childs := ToolFar(num, pid)
	var chridnum []gform.Data
	if childs != nil {
		var number int = 1
		var total int = len(childs)
		for _, v := range childs {
			j := ""
			k := ""
			if number == total {
				j += "└"
				k = ""
				if itemprefix != "" {
					k = "&nbsp;"
				}

			} else {
				j += "├"
				k = ""
				if itemprefix != "" {
					k = "│"
				}
			}
			spacer := ""
			if itemprefix != "" {
				spacer = itemprefix + j
			}
			v["spacer"] = spacer
			v["children"] = GetTreeArray(num, v["id"].(int64), itemprefix+k+"&nbsp;")
			chridnum = append(chridnum, v)
			number++
		}
	}
	return chridnum
}

func GetRuleTreeArray(num []gform.Data, pid int64, itemprefix string) []gform.Data {
	childs := ToolFar(num, pid)
	var chridnum []gform.Data
	if childs != nil {
		var number int = 1
		var total int = len(childs)
		for _, v := range childs {
			j := ""
			k := ""
			if number == total {
				j += "└"
				k = ""
				if itemprefix != "" {
					k = "&nbsp;"
				}

			} else {
				j += "├"
				k = ""
				if itemprefix != "" {
					k = "│"
				}
			}
			spacer := ""
			if itemprefix != "" {
				spacer = itemprefix + j
			}
			v["spacer"] = spacer
			v["children"] = GetTreeArray(num, v["id"].(int64), itemprefix+k+"&nbsp;")
			chridnum = append(chridnum, v)
			number++
		}
	}
	return chridnum
}

func GetRuleTreeArrayByPack(num []gform.Data, pid int64) []gform.Data {
	childs := ToolFar(num, pid)
	var chridnum []gform.Data
	if childs != nil {
		for _, v := range childs {
			newdata := GetRuleTreeArrayByPack(num, v["id"].(int64))
			if newdata != nil {
				v["children"] = GetRuleTreeArrayByPack(num, v["id"].(int64))
			}
			chridnum = append(chridnum, v)
		}
	}
	return chridnum
}

func GetMenuChildrenArray(pdata []gform.Data, parent_id int64, pid_file string) []gform.Data {
	var returnList []gform.Data
	for _, v := range pdata {
		if v[pid_file].(int64) == parent_id {
			children := GetMenuChildrenArray(pdata, v["id"].(int64), pid_file)
			if children != nil {
				v["children"] = children
			}
			returnList = append(returnList, v)
		}
	}
	return returnList
}

func GetMenuChildrenArraylist(pdata []gform.Data, parent_id int64) []gform.Data {
	var returnList []gform.Data
	for _, v := range pdata {
		if v["pid"].(int64) == parent_id {
			children := GetMenuChildrenArraylist(pdata, v["value"].(int64))
			if children != nil {
				v["children"] = children
			}
			returnList = append(returnList, v)
		}
	}
	return returnList
}

func ToolFar(data []gform.Data, pid int64) []gform.Data {
	var mapString []gform.Data
	for _, v := range data {
		if v["pid"].(int64) == pid {
			mapString = append(mapString, v)
		}
	}
	return mapString
}

func GetTreeList_txt(data []gform.Data, field string) []gform.Data {
	var midleArr []gform.Data
	for _, v := range data {
		var children []gform.Data
		if _, ok := v["children"]; ok {
			children = v["children"].([]gform.Data)
		} else {
			children = make([]gform.Data, 0)
		}
		delete(v, "children")
		v[field+"_txt"] = v["spacer"].(string) + " " + v[field+""].(string)
		if len(children) > 0 {
			v["haschild"] = 1
		} else {
			v["haschild"] = 0
		}
		if _, ok := v["id"]; ok {
			midleArr = append(midleArr, v)
		}
		if len(children) > 0 {
			newarr := GetTreeList_txt(children, field)
			midleArr = ArrayMerge_x(midleArr, newarr)
		}
	}
	return midleArr
}

func ArrayMerge_x(ss ...[]gform.Data) []gform.Data {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]gform.Data, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}
