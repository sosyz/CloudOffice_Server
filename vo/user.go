package vo

type LoginVo struct {
	Code string `form:"code" json:"code" uri:"code" xml:"code" binding:"required"`
	From string `form:"from" json:"from" uri:"from" xml:"from" binding:"required"`
}

type Code2SessionVo struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type TokenVo struct {
	Openid  string `form:"openid" json:"openid" uri:"openid" xml:"openid" binding:"required"`
	Session string `form:"session" json:"session" uri:"session" xml:"session" binding:"required"`
}

type UserInfo struct {
	Name    string `form:"name" json:"name" uri:"name" xml:"name" binding:"required"`
	Phone   string `form:"phone" json:"phone" uri:"phone" xml:"phone" binding:"required"`
	Group   uint   `form:"group" json:"group" uri:"group" xml:"group" binding:"required"`
	Address string `form:"address" json:"address" uri:"address" xml:"address" binding:"required"`
}
