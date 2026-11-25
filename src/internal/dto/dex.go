package dto

type GetDexResponse struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Color string `json:"color,omitempty"`
}

type CreateDexRequestBody struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

type GetDexsRequestBody RequestBody

type GetDexsResponse []GetDexResponse
