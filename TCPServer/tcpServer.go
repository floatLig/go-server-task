package TCPServer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	config "shopee.com/zeliang-entry-task/const"

	"shopee.com/zeliang-entry-task/rpc"

	"shopee.com/zeliang-entry-task/model"

	"shopee.com/zeliang-entry-task/redis"

	"shopee.com/zeliang-entry-task/mysql"

	_ "net/http/pprof"
)

func TCPServerMain() {

	//runtime.SetBlockProfileRate(1)
	//go func() {
	//	fmt.Println("启动监听")
	//	err := http.ListenAndServe("localhost:10001", nil)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()

	err := mysql.InitMySql()
	if err != nil {
		log.Println("mysql 初始化失败", err)
	}
	//insert()
	err = redis.InitRedis()
	if err != nil {
		fmt.Println("redis 初始化失败")
	}

	listen, err := net.Listen("tcp", ":3030")
	if err != nil {
		panic(err)
	}
	log.Println("listen to 3030")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("conn err:", err)
		} else {
			go HandleConn(conn)
		}
	}

}

func HandleConn(conn net.Conn) {
	for {
		data := rpc.RpcRead(conn)

		jsonType := model.JsonType{}
		_ = json.Unmarshal(data, &jsonType)

		var obj interface{} = nil
		switch jsonType.Type {
		case "login":
			obj = doLogin(data)
		case "update":
			obj = doUpdate(data)
		case "user":
			obj = doUser(data)
		default:
			fmt.Println("接受数据错误：", string(data))
		}

		data, _ = json.Marshal(obj)
		rpc.RpcSend(conn, &data)

	}
}

func doLogin(data []byte) *model.LoginRes {
	loginReq := model.LoginReq{}
	err := json.Unmarshal(data, &loginReq)
	if err != nil {
		fmt.Println("xx==> doLogin json unmarshal fail, err:", err.Error())
	}

	res := &model.LoginRes{}
	res.Success = true
	var userinfo *model.UserInfo

	// 从 redis 中取
	userKey := config.UserKeyPrefix + loginReq.UserName
	userinfo, err = redis.GetUser(userKey)
	//if err != nil {
	//	log.Println("get user from redis error, err:", err)
	//}

	// 如果redis找不到，从 mysql 中取
	if userinfo == nil {
		userinfo, err = mysql.SelectUser(loginReq.UserName)
		if userinfo == nil || err != nil {
			res.Success = false
			res.Error = "user no exist"
			return res
		}
	}

	newToken := loginReq.UserName + config.TokenSuffix
	_ = redis.SaveUserByToken(newToken, userinfo)

	md5String := GetMd5String(loginReq.Password)
	if userinfo.Password != md5String {
		res.Error = "user's password is not correct"
		res.Success = false
		return res
	}

	res.Token = newToken
	return res
}

func doUpdate(data []byte) *model.UserRes {
	updateReq := model.UpdateReq{}
	err := json.Unmarshal(data, &updateReq)
	if err != nil {
		fmt.Println("xx==> json unmarshal UpdateReq[obj] fail, error:", err.Error())
		return nil
	}

	updateRes := &model.UserRes{}
	updateRes.Success = true
	updateRes.Username = updateReq.Username
	updateRes.Nickname = updateReq.Nickname

	userInfo := updateReq.UpdateReqToUserInfo()
	if updateReq.Picture == "" {
		err = mysql.UpdateUserNickname(userInfo)
	} else {
		err = mysql.UpdateUserNicknameAndPicture(userInfo)
	}

	if err != nil {
		fmt.Println("xx==> update error")
		updateRes.Success = false
		updateRes.Error = err.Error()
		return updateRes
	}
	userKey := config.UserKeyPrefix + userInfo.Username
	redis.Del(userKey)

	return updateRes
}

func doUser(data []byte) *model.UserRes {
	res := &model.UserRes{}
	res.Success = true

	req := &model.UserReq{}
	err := json.Unmarshal(data, req)
	if err != nil {
		res.Error = err.Error()
		res.Success = false
		return res
	}
	if req.Token == "" {
		res.Error = "please login first"
		res.Success = false
		return res
	}

	userKey := redis.GetUserKeyByToken(req.Token)
	userinfo, err := redis.GetUser(userKey)

	if userinfo == nil {
		username := strings.Split(userKey, ":")[1]
		userinfo, _ = mysql.SelectUser(username)
		_ = redis.SaveUser(userinfo)
		if userinfo == nil {
			res.Error = "please login fisrt"
			res.Success = false
			return res
		}
	}

	res.Username = userinfo.Username
	res.Nickname = userinfo.Nickname
	res.Picture = userinfo.Picture
	return res
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func insert() {
	var user model.UserInfo
	var i int
	for i = 0; i < 1000; i++ {
		password := GetMd5String("1")
		user.Username = strconv.Itoa(i) + "u"
		user.Password = password
		user.Nickname = "nickname" + strconv.Itoa(i)
		user.Picture = "1u.jpeg"
		mysql.InsertTestData(user)
	}
}
