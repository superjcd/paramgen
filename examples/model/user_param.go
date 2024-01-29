package model

type UserForm struct {
	Name     string `form:"name"`
	Password string `form:"password"`
	Email    string `form:"email"`
	IsAdmin  int    `form:"is_admin"`
}

type UserJson struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  int    `json:"is_admin"`
}
