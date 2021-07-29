package mysql

import (
	"fmt"

	config "shopee.com/zeliang-entry-task/const"

	"shopee.com/zeliang-entry-task/model"
)

var UserChannel = make(chan model.UserInfo, 1024)

func SelectAllUserToChan() {
	query, err := DB.Query("select * from user")
	defer query.Close()
	if err != nil {
		fmt.Println("xx==> select * from user fail, error:", err.Error())
		return
	}
	for query.Next() {
		userInfo := model.UserInfo{}
		err = query.Scan(&userInfo.Username, &userInfo.Nickname, &userInfo.Password, &userInfo.Picture)
		if err != nil {
			fmt.Println("xx==> query.Scan error occurred, error:", err.Error())
			continue
		}
		UserChannel <- userInfo
	}
}

func SelectUser(username string) (*model.UserInfo, error) {
	userinfo := model.UserInfo{}
	query := fmt.Sprintf("select username, nickname, password, picture from %s where username=?", config.MysqlTable)
	err := DB.QueryRow(query, username).Scan(&userinfo.Username, &userinfo.Nickname, &userinfo.Password, &userinfo.Picture)
	if err != nil {
		return nil, err
	}
	return &userinfo, nil
}

func UpdateUserNickname(user model.UserInfo) error {
	prepare, err := DB.Prepare("update user set nickname = ? where username = ?")
	defer prepare.Close()

	if err != nil {
		fmt.Println("xx==> updateUserInfo fail, error:", err.Error())
		return err
	}
	_, err = prepare.Exec(user.Nickname, user.Username)
	if err != nil {
		fmt.Println("xx==> updateUserInfo fail, error:", err.Error())
		return err
	}
	return nil
}

func UpdateUserNicknameAndPicture(user model.UserInfo) error {
	prepare, err := DB.Prepare("update user set nickname = ?, picture = ? where username = ?")
	defer prepare.Close()

	if err != nil {
		fmt.Println("xx==> updateUserInfo fail, error:", err.Error())
		return err
	}
	_, err = prepare.Exec(user.Nickname, user.Picture, user.Username)
	if err != nil {
		fmt.Println("xx==> updateUserInfo fail, error:", err.Error())
		return err
	}
	return nil
}

func InsertTestData(user model.UserInfo) error {
	prepare, err := DB.Prepare("insert into user (username, nickname, password, picture) values (?, ?, ?, ?)")
	defer prepare.Close()

	if err != nil {
		fmt.Println("xx==> updateUserInfo fail, error:", err.Error())
		return err
	}
	_, err = prepare.Exec(user.Username, user.Nickname, user.Password, user.Picture)
	if err != nil {
		fmt.Println("xx==> updateUserInfo fail, error:", err.Error())
		return err
	}
	return nil
}
