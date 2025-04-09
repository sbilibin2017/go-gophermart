package types

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserGetParam struct {
	Login string `json:"login"`
}
