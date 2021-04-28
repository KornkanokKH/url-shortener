package generate

type ShortenerRequest struct {
	ShortCode    string `json:"short_code" valid:"required"`
	FullURL      string `json:"full_url" valid:"required,url"`
	ExpireDate   int64  `json:"expire_date" valid:"int"`
	NumberOfHits int    `json:"number_of_hits" valid:"required,int"`
}

type ShortenerResponse struct {
	ShortCode string `json:"short_code"`
	FullURL   string `json:"full_url"`
	ShortURL  string `json:"short_url"`
}
