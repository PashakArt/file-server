package service

type GetListParams struct {
	Limit       int
	FilterKey   string
	FilterVal   string
	TargetLogin string
}

type DocResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Mime    *string  `json:"mime,omitempty"`
	File    bool     `json:"file"`
	Public  bool     `json:"public"`
	Created string   `json:"created"`
	Grant   []string `json:"grant"`
}

type GetDocsResponse struct {
	Docs []DocResponse `json:"docs"`
}
