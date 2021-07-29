package model

import "html/template"

type UserInfo struct {
	Username string
	Nickname string
	Password string
	Picture  string
}

type UserVo struct {
	Username string
	Nickname string
	Picture  template.URL
}

type LoginVo struct {
	ErrorReason string
}
