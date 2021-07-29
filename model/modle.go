package model

type JsonType struct {
	Type string `json:"type"`
}

type LoginReq struct {
	JsonType
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserReq struct {
	JsonType
	Token string `json:"token"`
}

type UpdateReq struct {
	JsonType
	Token    string `json:"token"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Picture  string `json:"picture"`
}

type LoginRes struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Token   string `json:"token"`
}

type UserRes struct {
	Success  bool   `json:"success"`
	Error    string `json:"reason"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Picture  string `json:"picture"`
}

func (updateReq *UpdateReq) UpdateReqToUserInfo() UserInfo {
	return UserInfo{
		Username: updateReq.Username,
		Nickname: updateReq.Nickname,
		Picture:  updateReq.Picture,
	}
}
