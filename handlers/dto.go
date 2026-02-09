package handlers

type Link struct {
	Id          uint64 `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortUrl    string `json:"short_url"`
}

type LinkParams struct {
	OriginalUrl string `json:"original_url" binding:"required,url"`
	ShortName   string `json:"short_name,omitempty" binding:"omitempty,alphanum,min=6,max=50"`
}

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
