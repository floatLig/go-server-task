package HTTPServer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"

	config "shopee.com/zeliang-entry-task/const"

	"shopee.com/zeliang-entry-task/connectionpool"

	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
	"shopee.com/zeliang-entry-task/model"
	"shopee.com/zeliang-entry-task/rpc"
)

func connectToTcp() (net.Conn, error) {
	tcpAddress := fmt.Sprintf("%s:%s", config.TcpHost, config.TcpPort)
	socket, err := net.DialTimeout(config.TcpNetWork, tcpAddress, config.TcpTimeout)
	if err != nil {
		log.Println("connect failed, err :", err.Error())
		return nil, err
	}
	return socket, nil
}

func HttpServerMain() {
	connectionpool.InitConnectionPool(config.HTTPConnectionPoolSize, connectToTcp)
	defer connectionpool.Close()

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/update", update)
	http.HandleFunc("/info", info)

	httpAddress := fmt.Sprintf("%s:%s", config.HttpHost, config.HttpPort)
	err := http.ListenAndServe(httpAddress, nil)
	if err != nil {
		log.Fatalln("http 监听失败，", err.Error())
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles(config.LoginFile)
	if err != nil {
		log.Println("xx==> template parse fail, error:", err.Error())
		return
	}
	err = files.Execute(w, nil)
	if err != nil {
		log.Println("xx==> template execute fail, error:", err.Error())
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	tcpServerConn, err := connectionpool.Get()
	defer connectionpool.Put(tcpServerConn)
	if err != nil {
		fmt.Println("xx==> getting connection from connectionPool error:", err.Error())
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	loginVo := &model.LoginVo{}

	if username == "" || password == "" {
		loginVo.ErrorReason = "Please enter your user name or password"
		err := returnLoginPage(w, loginVo)
		if err != nil {
			log.Println("返回主页面失败")
		}
		return
	}

	loginReq := &model.LoginReq{
		JsonType: model.JsonType{Type: "login"},
		UserName: username,
		Password: password,
	}
	data, err := json.Marshal(loginReq)
	if err != nil {
		fmt.Println("xx==> json Marshal loginReq(obj) fail, error:", err.Error())
		return
	}
	rpc.RpcSend(tcpServerConn, &data)

	data = rpc.RpcRead(tcpServerConn)
	loginRes := &model.LoginRes{}
	err = json.Unmarshal(data, loginRes)
	if err != nil {
		fmt.Println("xx==> json UnMarshal loginRes(obj) fail, error:", err.Error())
		return
	}

	if !loginRes.Success {
		loginVo.ErrorReason = loginRes.Error
		err := returnLoginPage(w, loginVo)
		if err != nil {
			log.Println("返回登陆页面失败")
		}
		return
	}

	token := loginRes.Token
	c := http.Cookie{
		Name:     config.LoginToken,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/info", 302)
}

func update(w http.ResponseWriter, r *http.Request) {
	tcpServerConn, err := connectionpool.Get()
	defer connectionpool.Put(tcpServerConn)
	if err != nil {
		log.Println("xx==> getting connection from connectionPool error:", err.Error())
	}

	username := r.PostFormValue("username")
	nickname := r.PostFormValue("nickname")
	file, header, err := r.FormFile("picture")
	cookie, err := r.Cookie(config.LoginToken)
	if cookie == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	token := cookie.Value

	var pictureName string
	imgStorePath := config.ImgStorePath
	imgSuffix := []string{"png", "jpeg", "jpg", "gif"}

	createPath(imgStorePath)
	if header == nil || file == nil {
		pictureName = ""
	} else {
		// 将用户的照片名简化
		strs := strings.Split(header.Filename, ".")
		pictureName = username + "." + strs[len(strs)-1]
	}

	if pictureName != "" && judgeIsImg(pictureName, imgSuffix) {
		delOldPicture(username, imgSuffix, imgStorePath)
		saveNewPicture(pictureName, imgStorePath, file)
	}

	updateReq := model.UpdateReq{
		JsonType: model.JsonType{
			Type: "update",
		},
		Token:    token,
		Username: username,
		Nickname: nickname,
		Picture:  pictureName,
	}
	data, err := json.Marshal(updateReq)
	if err != nil {
		fmt.Println("xx==> json Marshal updateReq(obj) fail, error:", err.Error())
	}
	rpc.RpcSend(tcpServerConn, &data)

	data = rpc.RpcRead(tcpServerConn)
	updateRes := model.UserRes{}
	err = json.Unmarshal(data, &updateRes)
	if err != nil || !updateRes.Success {
		fmt.Println("xx==> json UnMarshal updateRes(obj) fail, error:", err.Error())
		http.Redirect(w, r, "/", 302)
		return
	}

	http.Redirect(w, r, "/info", 302)
}

// 图片 转成 base64
func returnImg(picture string) template.URL {
	var base64Str string

	picturePath := "./img/" + picture
	srcByte, err := ioutil.ReadFile(picturePath)
	if err != nil {
		log.Println("xx==> read local picture fail, error:", err.Error())
		return ""
	} else {
		base64Str = "data:image/" + picture + ";base64," + base64.StdEncoding.EncodeToString(srcByte)
	}
	return template.URL(base64Str)
}

func saveNewPicture(pictureName string, imgStorePath string, file multipart.File) {
	imgPath := imgStorePath + "/" + pictureName
	create, err := os.Create(imgPath)
	if err != nil {
		fmt.Println("xx==> create picture error, error:", err.Error())
		return
	}
	_, err = io.Copy(create, file)
	if err != nil {
		fmt.Println("xx==> copy picture error, error:", err.Error())
	}
	// todo all resource close
	_ = create.Close()
}

func delOldPicture(username string, imgSuffix []string, imgStorePath string) {
	for _, suffix := range imgSuffix {
		fileName := username + "." + suffix
		filePath := imgStorePath + "/" + fileName

		_, err := os.Stat(filePath) //os.Stat获取文件信息
		if err != nil {
			if os.IsExist(err) {
				_ = os.Remove(filePath)
			}
			continue
		}
		_ = os.Remove(filePath)
	}
}

func judgeIsImg(fileName string, imgSuffix []string) bool {
	if fileName == "" {
		return false
	}

	split := strings.Split(fileName, ".")
	fileSuffix := split[len(split)-1]

	res := false
	for _, i := range imgSuffix {
		if i == fileSuffix {
			res = true
			break
		}
	}
	return res
}

func createPath(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0777)
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func returnLoginPage(w http.ResponseWriter, vo *model.LoginVo) error {
	files, err := template.ParseFiles(config.LoginFile)
	if err != nil {
		return err
	}
	err = files.Execute(w, vo)
	if err != nil {
		return err
	}

	return nil
}

func returnUserPage(w http.ResponseWriter, vo *model.UserVo) error {
	// 返回用户信息
	//marshal, _ := json.Marshal(vo)
	//w.Write(marshal)
	//return nil

	files, err := template.ParseFiles(config.UserFile)
	if err != nil {
		fmt.Println("返回用户主页 打开失败")
		return err
	}
	err = files.Execute(w, vo)
	if err != nil {
		fmt.Println("返回用户主页 执行失败")
		return err
	}
	return nil
}

func info(w http.ResponseWriter, r *http.Request) {
	tcpServerConn, err := connectionpool.Get()
	defer connectionpool.Put(tcpServerConn)
	if err != nil {
		fmt.Println("xx==> getting connection from connectionPool error:", err.Error())
		return
	}
	cookie, _ := r.Cookie("token")
	if cookie == nil {
		vo := &model.LoginVo{
			ErrorReason: "please login first",
		}
		_ = returnLoginPage(w, vo)
		return
	}
	token := cookie.Value

	req := &model.UserReq{
		JsonType: model.JsonType{Type: "user"},
		Token:    token,
	}
	data, _ := json.Marshal(req)
	rpc.RpcSend(tcpServerConn, &data)

	data = rpc.RpcRead(tcpServerConn)

	res := model.UserRes{}
	err = json.Unmarshal(data, &res)
	if err != nil || !res.Success || res.Error != "" {
		vo := &model.LoginVo{
			ErrorReason: "please login first",
		}
		_ = returnLoginPage(w, vo)
		return
	}

	vo := &model.UserVo{
		Username: res.Username,
		Nickname: res.Nickname,
		Picture:  returnImg(res.Picture),
	}
	_ = returnUserPage(w, vo)
}
