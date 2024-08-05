package main

type Task struct {
	Id     int `json:"id"`
	Value  string `json:"value"`
	Done   bool   `json:"done"`
	UserId int    `json:"userId"`
}

type User struct {
	Id       int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
