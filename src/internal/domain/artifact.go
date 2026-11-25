package domain

type Extra map[string]string

type Artifact struct {
	ID        string
	Name      string
	Image     string
	Reference string
	Dex       Dex
	Sites     []Site
	Extra     Extra
}
