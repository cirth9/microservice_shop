package forms

type PasswordLoginForm struct {
	Mobile        string `form:"mobile" json:"mobile" binding:"required"`
	Password      string `form:"password" json:"password" binding:"required,min=3,max=20"`
	CaptchaNumber string `form:"captcha_number" json:"captcha_number" binding:"required"`
	CaptchaId     string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}

type UserRegister struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	CheckCode int    `form:"check_code" json:"check_code" binding:"required"`
}
