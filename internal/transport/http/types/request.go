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
	Name   string   `json:"name" validate:"required,min=1,max=255"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime" validate:"omitempty,max=100"`
	Grant  []string `json:"grant" validate:"omitempty,dive,required,alphanum"`
}
