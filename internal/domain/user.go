package domain

import "time"

// 领域对象，DDD中的 entity，BO
type User struct {
	Id       int64
	Email    string
	Password string
	Birthday string
	NickName string
	Profile  string
	Phone    string
	Ctime    time.Time
}
