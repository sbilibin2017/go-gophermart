package domain

type User struct {
	Login    string
	Password string
}

type UserGetParam struct {
	Login string
}
