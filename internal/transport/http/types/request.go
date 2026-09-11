package types

type LoginBody struct {
	Login    string `json:"login"`
	Password string `json:"pswd"`
}

type RegisterBody struct {
	LoginBody
	Token string `json:"token"`
}

type DocMeta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}
