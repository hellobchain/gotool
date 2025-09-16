package gvalid

import (
	"errors"
	"regexp"
	"strings"
)

// 常见正则预编译
var (
	regMobile = regexp.MustCompile(`^1[3-9]\d{9}$`)
	regEmail  = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	regIDCard = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
	regURL    = regexp.MustCompile(`^https?://[^\s/$.?#].\S*$`)
)

// ==== 单函数校验 ====
func IsMobile(s string) bool       { return regMobile.MatchString(s) }
func IsEmail(s string) bool        { return regEmail.MatchString(s) }
func IsIDCard(s string) bool       { return regIDCard.MatchString(s) }
func IsURL(s string) bool          { return regURL.MatchString(s) }
func IsDigits(s string) bool       { return regexp.MustCompile(`^\d+$`).MatchString(s) }
func IsLetters(s string) bool      { return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(s) }
func IsAlphanumeric(s string) bool { return regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(s) }

// NotBlank 非空且不只空白
func NotBlank(s string) bool { return strings.TrimSpace(s) != "" }

// ==== 链式校验器 ====
type Validator struct {
	value interface{}
	err   error
}

func New(v interface{}) *Validator { return &Validator{value: v} }

func (v *Validator) NotBlank(msg ...string) *Validator {
	if v.err != nil {
		return v
	}
	s, ok := v.value.(string)
	if !ok || !NotBlank(s) {
		v.err = buildErr(msg, "value must not be blank")
	}
	return v
}

func (v *Validator) Mobile(msg ...string) *Validator {
	if v.err != nil {
		return v
	}
	s, _ := v.value.(string)
	if !IsMobile(s) {
		v.err = buildErr(msg, "invalid mobile format")
	}
	return v
}

func (v *Validator) Email(msg ...string) *Validator {
	if v.err != nil {
		return v
	}
	s, _ := v.value.(string)
	if !IsEmail(s) {
		v.err = buildErr(msg, "invalid email format")
	}
	return v
}

func (v *Validator) Check() error { return v.err }

// 辅助
func buildErr(msg []string, def string) error {
	if len(msg) > 0 {
		return errors.New(msg[0])
	}
	return errors.New(def)
}
