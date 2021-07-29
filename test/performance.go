package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var cnt = 0

func main() {
	msgChan := make(chan string, 3000)
	fmt.Println("开始运行")
	for i := 0; i < 100; i++ {
		go makeApiRequests(2, msgChan)
	}
	for s := range msgChan {
		fmt.Println(s)
	}
}

func makeApiRequests(i int, msgChan chan string) {
	cnt++
	msg := fmt.Sprint("【", cnt, "】")

	username := strconv.Itoa(i) + "u"
	password := strconv.Itoa(i) + "p"
	reqPar := fmt.Sprintf("username=%s&password=%s", username, password)

	url := "http://localhost:8888/login"
	method := "POST"

	reqBody := strings.NewReader(reqPar)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		msgChan <- msg + err.Error()
		return
	}
	defer res.Body.Close()

	//msgChan <- msg + reqPar + "success"

	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(string(body))

	//================

	//if loginRespByt, err := loginCmd.Output(); err != nil {
	//	errChan <- fmt.Errorf("login curl failed")
	//	return
	//} else {
	//	if strings.Contains(string(loginRespByt), "修改") {
	//		outputChan <- fmt.Sprintf("%s Sucess", username)
	//	} else {
	//		errChan <- fmt.Errorf("login %s Failed", username)
	//		return
	//	}
	//}

	// update
	//updateUsername := reqUserName
	//updateName := fmt.Sprintf("name=%s", "performance test")
	//updatePassword := fmt.Sprintf("pass_word=%s", "12345678")
	//updatePic := fmt.Sprintf("pic_profile=%s", "@/Users/hao.wu/GolandProjects/entrytask/http/template/file/test2_729365000.png")
	//updateCmd := exec.Command("curl", "-H", "Cookie:entry_task_session=0z4UwRJrWp_mX7RnKYSiqHwkspjraX5RlzECP5vIcfY=", "-F", updateUsername, "-F", updateName, "-F", updatePassword, "-F", updatePic, "http://127.0.0.1:8080/update")
	//if updateByt, err := updateCmd.Output(); err != nil {
	//	errChan <- fmt.Errorf("update curl failed")
	//	return
	//} else {
	//	if strings.Contains(string(updateByt), "修改") {
	//		outputChan <- fmt.Sprintf("%s Success", username)
	//	} else {
	//		errChan <- fmt.Errorf("update %s Failed", username)
	//		return
	//	}
	//}
	//
	//// logout
	//logoutUsername := fmt.Sprintf("{\"%s\":\"%s\"}", "user_name", username)
	//logoutCmd := exec.Command("curl", "-H", "Cookie:entry_task_session=0z4UwRJrWp_mX7RnKYSiqHwkspjraX5RlzECP5vIcfY=", "-H", "Content-Type: application/json", "-d", logoutUsername, "http://127.0.0.1:8080/logout")
	//if logoutByt, err := logoutCmd.Output(); err != nil {
	//	errChan <- fmt.Errorf("logout curl failed")
	//	return
	//} else {
	//	if strings.Contains(string(logoutByt), "Hello") {
	//		outputChan <- fmt.Sprintf("%s Success", username)
	//	} else {
	//		errChan <- fmt.Errorf("logout %s Failed", username)
	//	}
	//}
}
