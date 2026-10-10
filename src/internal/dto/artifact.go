package dto

type Measurements struct {
	Chest float32 `json:"chest,omitempty"`
	Hips  float32 `json:"hips,omitempty"`
	Waist float32 `json:"waist,omitempty"`
}

type Site struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GetArtifactResponse struct {
	ID           string            `json:"id" validate:"required"`
	Name         string            `json:"name" validate:"required"`
	Gener        string            `json:"gender" validate:"required,oneof=male female other"`
	Age          int               `json:"age,omitempty"`
	Height       float32           `json:"height,omitempty"`
	Birthdate    string            `json:"birthdate,omitempty"`
	Image        string            `json:"image,omitempty"`
	Measurements Measurements      `json:"measurements"`
	Element      string            `json:"element,omitempty"`
	Sites        []Site            `json:"sites,omitempty"`
	Domain       *Site             `json:"domain,omitempty"`
	Extra        map[string]string `json:"extra,omitempty"`
}

type CreateArtifactRequestBody struct {
	Name         string            `json:"name"`
	Gener        string            `json:"gender" validate:"oneof=male female other"`
	Height       float32           `json:"height,omitempty"`
	Birthdate    string            `json:"birthdate,omitempty" validate:"datetime=2006-01-02"`
	Image        string            `json:"image,omitempty"`
	Measurements Measurements      `json:"measurements"`
	Element      string            `json:"element,omitempty"`
	Sites        []Site            `json:"sites,omitempty" validate:"dive"`
	Extra        map[string]string `json:"extra,omitempty"`
}

type GetArtifactsRequestQuery RequestQuery

type GetArtifactsResponse []GetArtifactResponse
