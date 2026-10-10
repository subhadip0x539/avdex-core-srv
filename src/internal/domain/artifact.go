package domain

import "time"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Measurements struct {
	Chest float32
	Hips  float32
	Waist float32
}

type Extra map[string]string

type Site struct {
	ID   string
	Name string
	URL  string
}

type Artifact struct {
	ID           string
	Name         string
	Gender       Gender
	Age          int
	Height       float32
	Birthdate    time.Time
	Image        string
	Measurements Measurements
	Element      string
	Sites        []Site
	Domain       *Site
	Extra        Extra
}
