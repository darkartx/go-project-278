package handlers

type Link struct {
	Id          int64  `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortUrl    string `json:"short_url"`
}

type LinkParams struct {
	OriginalUrl string `json:"original_url"`
	ShortName   string `json:"short_name,omitempty"`
}
