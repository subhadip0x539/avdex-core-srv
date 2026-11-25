package dto

type Sort struct {
	By    string `json:"by"`
	Order string `json:"order"`
}

type Filter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Values   []any  `json:"values"`
}

type RequestBody struct {
	Sort    Sort     `json:"sort"`
	Filters []Filter `filters:"filters"`
}
