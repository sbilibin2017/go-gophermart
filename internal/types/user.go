package types

type UserID struct {
	UserID int64 `json:"user_id"`
}

type User struct {
	*UserID
}

func NewUser(userID int64) *User {
	return &User{
		UserID: &UserID{UserID: userID},
	}
}
