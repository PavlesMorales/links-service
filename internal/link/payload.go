package link

type LinkCreateRq struct {
	Url string `json:"url" validate:"required,url"`
}

type LinkUpdateRq struct {
	Url  string `json:"url" validate:"url"`
	Hash string `json:"hash"`
}

type GetAllLinksResponse struct {
	Links []Link `json:"links"`
	Count int64  `json:"count"`
}
