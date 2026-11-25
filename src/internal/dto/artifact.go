package dto

type Site struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GetArtifactResponse struct {
	ID        string            `json:"id" validate:"required"`
	Name      string            `json:"name" validate:"required"`
	Image     string            `json:"image,omitempty"`
	Reference string            `json:"reference,omitempty"`
	Sites     []Site            `json:"sites,omitempty"`
	DexID     string            `json:"dex_id,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

type CreateArtifactRequestBody struct {
	Name      string            `json:"name"`
	Image     string            `json:"image,omitempty"`
	Reference string            `json:"reference,omitempty"`
	DexID     string            `json:"dex_id,omitempty"`
	Sites     []Site            `json:"sites,omitempty" validate:"dive"`
	Extra     map[string]string `json:"extra,omitempty"`
}

type GetArtifactsRequestBody RequestBody

type GetArtifactsResponse []GetArtifactResponse
