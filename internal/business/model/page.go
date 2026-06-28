package model

type PageData struct {
	Title       string
	Description string
	Page        string
	Data        interface{}
	Error       ErrorData
}
