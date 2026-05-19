package gf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gofly/global"
	"gofly/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpGet(url_text string, data map[string]interface{}) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	for k, v := range data {
		paras.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("request token err :" + err.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New("request token response json parse err :" + err.Error())
	} else {
		return jMap, nil
	}

}

func HttpPost(url_text string, urldata map[string]interface{}, postdata map[string]interface{}, contentType string) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	for k, v := range urldata {
		paras.Set(k, v.(string))
	}
	u.RawQuery = paras.Encode()
	jsonData := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(jsonData)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(postdata); err != nil {
		return nil, errors.New("请求错误 :" + err.Error())
	}
	body := bytes.NewBufferString(string(jsonData.Bytes()))
	resp, erro := http.Post(u.String(), contentType, body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if erro != nil {
		return nil, errors.New("请求错误 :" + erro.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}

}

func HttpPost_c(url_text string, urldata map[string]interface{}, postdata map[string]interface{}, contentType string) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	for k, v := range urldata {
		paras.Set(k, v.(string))
	}
	u.RawQuery = paras.Encode()
	jsonStr, _ := json.Marshal(postdata)
	body := bytes.NewBuffer([]byte(jsonStr))
	resp, erro := http.Post(u.String(), contentType, body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if erro != nil {
		return nil, errors.New("请求错误 :" + erro.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}

}

type Response struct {
	Code      int         `json:"code"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
}

func ServerError(c *gin.Context, err interface{}) {
	conf := global.App.Config
	msg := "内部服务器错误"
	if os.Getenv(gin.EnvGinMode) != gin.ReleaseMode && reflect.TypeOf(err).Name() == "string" {
		msg = err.(string)
	} else {
		if conf.App.Env != "pro" && os.Getenv(gin.EnvGinMode) != gin.ReleaseMode {
			if _, ok := err.(error); ok {
				msg = err.(error).Error()
			}
		} else {
			str := fmt.Sprintf("内部服务器错误： %s\n", err.(error).Error())
			global.App.Log.Error(str)
		}
	}
	if res := strings.Contains(msg, "token is expired by"); res {
		c.JSON(200, Response{
			401,
			http.StatusInternalServerError,
			nil,
			msg,
		})
	} else if res := strings.Contains(msg, "invalid memory address or nil pointer dereference"); res {
		model.MyInit(3)
		c.JSON(http.StatusInternalServerError, Response{1,
			http.StatusInternalServerError,
			"可能是数据库链接失败，请查看数据库链接是否正常",
			msg + "，可能是数据库链接失败，请查看数据库配置及是否启动，再刷新试试！",
		})
	} else {
		c.JSON(http.StatusInternalServerError, Response{1,
			http.StatusInternalServerError,
			nil,
			msg,
		})
	}
	c.Abort()
}

func Get_x(url string) (string, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}
	return result.String(), nil
}

func Post_strdata(url string, data string, contentType string) (string, error) {
	if contentType == "" {
		contentType = "application/json"
	}
	payload := strings.NewReader(data)
	req, err := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", contentType)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		return "", error
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	return string(result), nil
}

// tool
func Get(url string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}

func Post(url string, data interface{}, contentType string) string {

	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result)
}
