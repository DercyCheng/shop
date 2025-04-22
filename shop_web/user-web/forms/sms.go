package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}
