package developer

import (
	"bufio"
	"fmt"
	"gofly/utils/gf"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gofly/utils/gform"
)

func CreatApicodeFile(model_name string, data gform.Data) {
	url := data["url"].(string)
	url_arr := strings.Split(url, `/`)
	methods := url_arr[len(url_arr)-1]
	filename := url_arr[len(url_arr)-2]
	model_path := strings.Split(url, filename)
	folder_path := model_path[0]
	packageName_path := strings.Split(folder_path, `/`)
	packageName := packageName_path[len(packageName_path)-2]
	file_path := filepath.Join("app/", folder_path)
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(file_path, os.ModePerm)
		}
	}
	filego_path := filepath.Join("app/", folder_path, filename+".go")
	if _, err := os.Stat(filego_path); err != nil {
		if !os.IsExist(err) {
			os.Create(filego_path)
			err := CopyFileContents(filepath.Join("resource/developer/codetpl/go/apicode.gos"), filego_path)
			if err != nil {
				panic(err)
			}
			ChangPackage(filego_path, packageName, filename, gf.InterfaceTostring(data["tablename"]))
			CheckIsAddController(model_name, model_name+"/"+packageName)
		}
	}
	if data["method"] == "get" {
		if data["getdata_type"] == "list" {
			CreatList(filego_path, methods, data["fields"].(string))
		} else if data["getdata_type"] == "detail" {
			CreatDetail(filego_path, methods, data["fields"].(string))
		}
	} else if data["method"] == "post" {
		CreatSave(filego_path, methods)
	} else if data["method"] == "delete" {
		CreatDel(filego_path, methods)
	}
}

func CreatList(filePath, methods, fields string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), "get_list") {
			datestr := strings.ReplaceAll(string(a), "get_list", gf.FirstUpper(methods))
			result += datestr + "\n"
		} else if strings.Contains(string(a), "{fields}") {
			datestr := strings.ReplaceAll(string(a), "{fields}", fields)
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func CreatDetail(filePath, methods, fields string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), "get_detail") {
			datestr := strings.ReplaceAll(string(a), "get_detail", gf.FirstUpper(methods))
			result += datestr + "\n"
		} else if strings.Contains(string(a), "{fields}") {
			datestr := strings.ReplaceAll(string(a), "{fields}", fields)
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func CreatSave(filePath, methods string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), "save(") {
			datestr := strings.ReplaceAll(string(a), "save(", gf.FirstUpper(methods)+"(")
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func CreatDel(filePath, methods string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), "del(") {
			datestr := strings.ReplaceAll(string(a), "del(", gf.FirstUpper(methods)+"(")
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func ChangPackage(filePath, packageName, filename, tablename string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), "packageName") {
			datestr := strings.ReplaceAll(string(a), "packageName", packageName)
			result += datestr + "\n"
		} else if strings.Contains(string(a), "Replace") {
			datestr := strings.ReplaceAll(string(a), "Replace", gf.FirstUpper(filename))
			result += datestr + "\n"
		} else if strings.Contains(string(a), "{tablename}") {
			datestr := strings.ReplaceAll(string(a), "{tablename}", tablename)
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func UnApicodeFile(data gform.Data) {
	url := data["url"].(string)
	url_arr := strings.Split(url, `/`)
	methods := url_arr[len(url_arr)-1]
	filename := url_arr[len(url_arr)-2]
	model_path := strings.Split(url, filename)
	folder_path := model_path[0]
	filego_path := filepath.Join("app/", folder_path, filename+".go")
	if _, err := os.Stat(filego_path); err == nil {
		if data["method"] == "get" {
			if data["getdata_type"] == "list" {
				UnList(filego_path, methods, gf.InterfaceTostring(data["fields"]))
			} else if data["getdata_type"] == "detail" {
				UnDetail(filego_path, methods, gf.InterfaceTostring(data["fields"]))
			}
		} else if data["method"] == "post" {
			UnSave(filego_path, methods)
		} else if data["method"] == "delete" {
			UnDel(filego_path, methods)
		}
	}
}

func UnList(filePath, methods, fields string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), gf.FirstUpper(methods)) {
			datestr := strings.ReplaceAll(string(a), gf.FirstUpper(methods), "get_list")
			result += datestr + "\n"
		} else if strings.Contains(string(a), fields) {
			datestr := strings.ReplaceAll(string(a), fields, "{fields}")
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func UnDetail(filePath, methods, fields string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), gf.FirstUpper(methods)) {
			datestr := strings.ReplaceAll(string(a), gf.FirstUpper(methods), "get_detail")
			result += datestr + "\n"
		} else if strings.Contains(string(a), fields) {
			datestr := strings.ReplaceAll(string(a), fields, "{fields}")
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func UnSave(filePath, methods string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), gf.FirstUpper(methods)+"(") {
			datestr := strings.ReplaceAll(string(a), gf.FirstUpper(methods)+"(", "save(")
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func UnDel(filePath, methods string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), gf.FirstUpper(methods)+"(") {
			datestr := strings.ReplaceAll(string(a), gf.FirstUpper(methods)+"(", "del(")
			result += datestr + "\n"
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func RemoveModel(model_name string, data gform.Data) {
	url := data["url"].(string)
	url_arr := strings.Split(url, `/`)
	filename := url_arr[len(url_arr)-2]
	model_path := strings.Split(url, filename)
	folder_path := model_path[0]
	filego_path := filepath.Join("app/", folder_path, filename+".go")
	if _, err := os.Stat(filego_path); err == nil {
		os.Remove(filego_path)
		file_path := fmt.Sprintf("%s%s", "app/", folder_path)
		dir, _ := os.ReadDir(file_path)
		if len(dir) == 0 {
			os.RemoveAll(file_path)
			packageName_path := strings.Split(folder_path, `/`)
			packageName := packageName_path[len(packageName_path)-2]
			CheckApiRemoveController(model_name, model_name+"/"+packageName)
		}
	}
}
