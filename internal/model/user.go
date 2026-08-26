package model

import (
	"strings"
	"time"
)

// User 表示系统用户。
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	if u.Username == "" {
		return NewValidationError("username", "用户名不能为空")
	}
	if u.Email == "" {
		return NewValidationError("email", "邮箱不能为空")
	}
	return nil
}

type UserFilter struct {
	Keyword string
}

func (f UserFilter) Match(u *User) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(u.Username), k) &&
			!strings.Contains(strings.ToLower(u.Email), k) {
			return false
		}
	}
	return true
}
