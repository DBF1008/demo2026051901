package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"gofly/model"
	"gofly/utils/gf"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"gofly/utils/gform"

	"github.com/gin-gonic/gin"
)

func init() {
	gf.Register(&Index{}, reflect.TypeOf(Index{}).PkgPath())
}

type Index struct{}

func (api *Index) GetPost_api(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	Data, err := model.DB().Table("business_wxsys_officonfig").Where("businessID", id).Fields("Token,accountID,businessID").First()
	if err != nil {
		log.Println("获取账号信息失败!")
		return
	} else {
		if c.Request.Method == "GET" {
			signature := c.DefaultQuery("signature", "")
			timestamp := c.DefaultQuery("timestamp", "")
			nonce := c.DefaultQuery("nonce", "")
			echostr := c.DefaultQuery("echostr", "")
			ok := CheckSignature(signature, timestamp, nonce, Data["Token"].(string))
			if !ok {
				log.Println("微信公众号接入校验失败!")
				return
			}
			log.Println("微信公众号接入校验成功!")
			_, _ = c.Writer.WriteString(echostr)
		} else {
			log.Println("post请求!")
			openid := c.DefaultQuery("openid", "")
			log.Printf("Openid: %s\n", openid)
			postdata, err := c.GetRawData()
			if err != nil {
				log.Fatalln(err)
			}
			msgText, err := ReceiveCommonMsg(postdata)
			if err != nil {
				log.Fatalln(err)
			}
			if msgText.Event == "subscribe" {

				Onsubscribe(msgText, openid, Data)
			} else if msgText.Event == "unsubscribe" {
				model.DB().Table("business_wxsys_user").Where("openid", openid).Data(map[string]interface{}{"subscribe": 0}).Update()
			}
			log.Printf("[消息接收] - 收到消息, 消息类型为: %s, FromUserName: %s\n", msgText.Event, msgText.FromUserName)
		}
	}
}

func Onsubscribe(msgText WxReceiveCommonMsg, openid string, Data gform.Data) {
	log.Println("判断账号是否存在!")
	user, _ := model.DB().Table("business_wxsys_user").Where("openid", openid).Fields("id").First()
	if user == nil {
		log.Println("新增账号!")
		userid, err := model.DB().Table("business_wxsys_user").Data(map[string]interface{}{
			"openid":     openid,
			"accountID":  Data["accountID"],
			"businessID": Data["businessID"],
			"subscribe":  1,
			"avatar":     "resource/staticfile/avatar.png",
			"createtime": time.Now().Unix(),
		}).InsertGetId()
		model.DB().Table("business_wxsys_user").Data(map[string]interface{}{
			"nickname": fmt.Sprintf("U_%v", userid),
		}).Where("id", userid).Update()
		log.Printf("添加失败: %s\n", err)
	} else {
		model.DB().Table("business_wxsys_user").Where("id", user["id"]).Data(map[string]interface{}{"subscribe": 1}).Update()
	}
}

type WxReceiveCommonMsg struct {
	ToUserName   string
	FromUserName string
	Content      string
	CreateTime   int64
	MsgType      string
	MsgId        int64
	PicUrl       string
	MediaId      string
	Event        string
	EventKey     string
	MenuId       string
	Format       string
	Recognition  string
	ThumbMediaId string
}

var WxReceiveFunc func(msg WxReceiveCommonMsg) error

func ReceiveCommonMsg(msgData []byte) (WxReceiveCommonMsg, error) {
	fmt.Printf("received weixin msgData:\n%s\n", msgData)
	msg := WxReceiveCommonMsg{}
	err := xml.Unmarshal(msgData, &msg)
	if WxReceiveFunc == nil {
		return msg, err
	}
	err = WxReceiveFunc(msg)
	return msg, err
}

func CheckSignature(signature, timestamp, nonce, token string) bool {
	arr := []string{timestamp, nonce, token}
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(arr); i++ {
		b.WriteString(arr[i])
	}

	return Sha1(b.String()) == signature
}

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
