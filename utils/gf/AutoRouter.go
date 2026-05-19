package gf

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Route struct {
	path       string
	httpMethod string
	Method     reflect.Value
	Args       []reflect.Type
}

var Routes = []Route{}

func Register(controller interface{}, PkgPathstr string) bool {
	vbf := reflect.ValueOf(controller)
	if vbf.NumMethod() == 0 {
		return false
	}
	rootPkg := ""
	if strings.Contains(PkgPathstr, "/app") {
		PkgPath_arr := strings.Split(PkgPathstr, "/app")
		rootPkg = PkgPath_arr[len(PkgPath_arr)-1]
	}
	ctrlName := reflect.TypeOf(controller).String()
	// fmt.Println("ctrlName=", ctrlName)
	module := ctrlName
	if strings.Contains(ctrlName, ".") {
		module = ctrlName[strings.Index(ctrlName, ".")+1:]
	}
	// fmt.Println("module=", module)
	if module == "Index" {
		module = "/"
	} else {
		module = "/" + strings.ToLower(module) + "/"
	}
	v := reflect.ValueOf(controller)
	// fmt.Println(ctrlName)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		action := v.Type().Method(i).Name
		path := rootPkg + module + FirstLower(action)
		params := make([]reflect.Type, 0, v.NumMethod())
		httpMethod := "POST"
		if (strings.HasPrefix(action, "Get") && !strings.HasPrefix(action, "GetPost")) || action == "Index" {
			httpMethod = "GET"
		} else if strings.HasPrefix(action, "Del") || action == "Del" {
			httpMethod = "DELETE"
		} else if strings.HasPrefix(action, "Put") || action == "Put" {
			httpMethod = "PUT"
		}
		for j := 0; j < method.Type().NumIn(); j++ {
			params = append(params, method.Type().In(j))
		}
		// fmt.Println("params=", params)
		// fmt.Println("action=", action)
		route := Route{path: path, Method: method, Args: params, httpMethod: httpMethod}
		Routes = append(Routes, route)
		if strings.HasPrefix(action, "GetPost") {
			route := Route{path: path, Method: method, Args: params, httpMethod: "GET"}
			Routes = append(Routes, route)
		}
	}
	// fmt.Println("Routes=", Routes)
	return true
}

func Bind(e *gin.Engine) {
	for _, route := range Routes {
		if route.httpMethod == "GET" {
			e.GET(route.path, match(route.path, route))
		}
		if route.httpMethod == "POST" {
			e.POST(route.path, match(route.path, route))
		}
		if route.httpMethod == "DELETE" {
			e.DELETE(route.path, match(route.path, route))
		}
		if route.httpMethod == "PUT" {
			e.PUT(route.path, match(route.path, route))
		}
	}
}

func match(path string, route Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(path, "/")
		if len(fields) < 3 {
			return
		}
		if len(Routes) > 0 {
			arguments := make([]reflect.Value, 1)
			arguments[0] = reflect.ValueOf(c) // *gin.Context
			route.Method.Call(arguments)
		}
	}
}
