package link

type CreateLinkRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type CreateLinkResponse struct {
	Id int `json:"id"`
	Code string `json:"code"`
	Url string `json:"url"`
}

type FindLinkResponse struct {
	Id int `json:"id"`
	Code string `json:"code"`
	Url string `json:"url"`
	Clicks int `json:"clicks"`
}