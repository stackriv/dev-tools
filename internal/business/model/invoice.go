package model

type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
}

type InvoiceData struct {
	Number   string `json:"number"`
	Date     string `json:"date"`
	DueDate  string `json:"dueDate"`
	Currency string `json:"currency"`
	From     struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
	} `json:"from"`
	To struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	} `json:"to"`
	Items   []InvoiceItem `json:"items"`
	Notes   string        `json:"notes"`
	TaxRate float64       `json:"taxRate"`
}
