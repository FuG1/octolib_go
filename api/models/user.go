package models

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Role         int    `json:"role"`
	ID           int    `json:"id"`
	FullName     string `json:"full_name"`
	TicketNumber string `json:"ticket_number"`
	BirthDate    string `json:"birth_date"`
	Phone        string `json:"phone"`
	Education    string `json:"education"`
	HallID       int    `json:"hall_id"`
	Books        []Book `json:"books"`
}
