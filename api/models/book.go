package models

type Book struct {
	Title         string `json:"title"`
	Author        int    `json:"author_id"`
	Genre         int    `json:"genre_id"`
	Description   string `json:"description"`
	PublishedDate string `json:"published_date"`
	Popularity    int    `json:"popularity"`
}
