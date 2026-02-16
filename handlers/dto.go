package handlers

import "time"

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

type Visit struct {
	Id        uint64    `json:"id"`
	LinkId    uint64    `json:"link_id"`
	Ip        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
