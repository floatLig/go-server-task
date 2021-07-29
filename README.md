
### Login API



```go
type LoginReq struct {
	JsonType
	UserName string `json:"username"`
	Password string `json:"password"`
}
```