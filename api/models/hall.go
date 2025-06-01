package models

type Hall struct {
	ID             int    `json:"id"`
	LibraryName    string `json:"library_name"`
	HallName       string `json:"hall_name"`
	Specialization string `json:"specialization"`
	Capacity       int    `json:"capacity"`
	Readers        []User `json:"readers"`
}
