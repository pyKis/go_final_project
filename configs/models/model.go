package models

const DatePattern string = "20060102"

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskIdResponse struct {
	Id uint `json:"id"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Sign struct {
	Password string `json:"password"`
}

type AuthToken struct {
	Token string `json:"token"`
}