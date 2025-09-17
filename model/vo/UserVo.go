package vo

type UserVo struct {
	UserName string
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Age      int    `json:"age"`
}
