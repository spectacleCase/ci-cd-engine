package utils

import "github.com/go-playground/validator/v10"

// EmailVerify 校验邮箱
func EmailVerify(email string) bool {
	v := validator.New()
	err := v.Var(email, "required,email")
	return err == nil
}
