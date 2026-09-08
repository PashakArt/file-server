package types

type LoginBody struct {
	Login    string `json:"login"`
	Password string `json:"pswd"`
}

type RegisterBody struct {
	LoginBody
	Token string `json:"token"`
}
