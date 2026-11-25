package handler

import (
	"log/slog"
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
		slog.Error(err.Trace.Error())
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	sites := []dto.Site{}

	for _, s := range result.Sites {
		sites = append(sites, dto.Site{
			Name: s.Name,
			URL:  s.URL,
		})
	}

	artifact := dto.GetArtifactResponse{
		ID:        result.ID,
		Name:      result.Name,
		Image:     result.Image,
		Sites:     sites,
		Reference: result.Reference,
		DexID:     result.Dex.ID,
		Extra:     result.Extra,
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.GetArtifactResponse]{Body: artifact})
}

func (h *ArtifactHandler) GetArtifacts(ctx *gin.Context) {
	var body dto.GetArtifactsRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response[any]{Message: err.Error()})
		return
	}

	result, err := h.svc.GetArtifacts(ctx, body)
	if err != nil {
		slog.Error(err.Trace.Error())
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	artifacts := dto.GetArtifactsResponse{}

	for _, r := range result {
		sites := []dto.Site{}

		for _, s := range r.Sites {
			sites = append(sites, dto.Site{
				Name: s.Name,
				URL:  s.URL,
			})
		}

		artifacts = append(artifacts, dto.GetArtifactResponse{
			ID:        r.ID,
			Name:      r.Name,
			Image:     r.Image,
			Sites:     sites,
			Reference: r.Reference,
			DexID:     r.Dex.ID,
			Extra:     r.Extra,
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
		slog.Error(err.Trace.Error())
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
