package model

type Starter struct {
	StringData map[string]string
	IntData    map[string]int
	Error      ErrorData
}

type ErrorData struct {
	Code    string
	Message string
}
