package domain

type User struct {
	Login    string
	Password string
}

type UserToken struct {
	Access string
}
