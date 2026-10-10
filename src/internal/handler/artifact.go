package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/port"
)

type ArtifactHandler struct {
	svc port.ArtifactService
}

func NewArtifactHandler(service port.ArtifactService) *ArtifactHandler {
	return &ArtifactHandler{svc: service}
}

func (h *ArtifactHandler) GetArtifact(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Artifact id is required"})
		return
	}

	result, err := h.svc.GetArtifact(ctx, id)
	if err != nil {
		ctx.JSON(err.Code, dto.Response[any]{Message: err.Message})
		return
	}

	sites := []dto.Site{}

	for _, s := range result.Sites {
		sites = append(sites, dto.Site{
			ID:   s.ID,
			Name: s.Name,
			URL:  s.URL,
		})
	}

	artifact := dto.GetArtifactResponse{
		ID:           result.ID,
		Name:         result.Name,
		Gener:        string(result.Gender),
		Height:       result.Height,
		Birthdate:    result.Birthdate.Format("2006-01-02"),
		Age:          result.Age,
		Measurements: dto.Measurements(result.Measurements),
		Image:        result.Image,
		Element:      result.Element,
		Sites:        sites,
		Domain:       (*dto.Site)(result.Domain),
		Extra:        result.Extra,
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.GetArtifactResponse]{Body: artifact})
}

func (h *ArtifactHandler) GetArtifacts(ctx *gin.Context) {
	var query dto.GetArtifactsRequestQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response[any]{Message: err.Error()})
		return
	}

	result, err := h.svc.GetArtifacts(ctx, query)
	if err != nil {
		ctx.JSON(err.Code, dto.Response[any]{Message: err.Message})
		return
	}

	artifacts := dto.GetArtifactsResponse{}

	for _, r := range result {
		sites := []dto.Site{}

		for _, s := range r.Sites {
			sites = append(sites, dto.Site{
				ID:   s.ID,
				Name: s.Name,
				URL:  s.URL,
			})
		}

		artifacts = append(artifacts, dto.GetArtifactResponse{
			ID:           r.ID,
			Name:         r.Name,
			Gener:        string(r.Gender),
			Height:       r.Height,
			Birthdate:    r.Birthdate.Format("2006-01-02"),
			Age:          r.Age,
			Measurements: dto.Measurements(r.Measurements),
			Image:        r.Image,
			Element:      r.Element,
			Sites:        sites,
			Domain:       (*dto.Site)(r.Domain),
			Extra:        r.Extra,
		})
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.GetArtifactsResponse]{Body: artifacts})
}

func (h *ArtifactHandler) CreateArtifact(ctx *gin.Context) {
	var body dto.CreateArtifactRequestBody
	if ver := ctx.ShouldBindJSON(&body); ver != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response[any]{Message: ver.Error()})
		return
	}

	_, err := h.svc.CreateArtifact(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response[any]{
		Message: "Artifact created successfully",
	})
}

func (h *ArtifactHandler) CreateDomain(ctx *gin.Context) {
	var body dto.CreateArtifactRequestBody
	if ver := ctx.ShouldBindJSON(&body); ver != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response[any]{Message: ver.Error()})
		return
	}

	_, err := h.svc.CreateArtifact(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response[any]{
		Message: "Artifact created successfully",
	})
}
