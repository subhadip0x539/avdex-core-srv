package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/port"
)

type DexHandler struct {
	svc port.DexService
}

func NewDexHandler(service port.DexService) *DexHandler {
	return &DexHandler{svc: service}
}

func (h *DexHandler) GetDex(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dex id is required"})
		return
	}

	result, err := h.svc.GetDex(ctx, id)
	if err != nil {
		slog.Error(err.Trace.Error())
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	dex := dto.GetDexResponse{
		ID:    result.ID,
		Name:  result.Name,
		Color: result.Color,
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.GetDexResponse]{Body: dex})
}

func (h *DexHandler) GetDexs(ctx *gin.Context) {
	var body dto.GetDexsRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response[any]{Message: err.Error()})
		return
	}

	result, err := h.svc.GetDexs(ctx, body)
	if err != nil {
		slog.Error(err.Trace.Error())
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	dexs := []dto.GetDexResponse{}

	for _, r := range result {
		dexs = append(dexs, dto.GetDexResponse{
			ID:    r.ID,
			Name:  r.Name,
			Color: r.Color,
		})
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.GetDexsResponse]{Body: dexs})
}

func (h *DexHandler) CreateDex(ctx *gin.Context) {
	var body dto.CreateDexRequestBody
	if ver := ctx.ShouldBindJSON(&body); ver != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response[any]{Message: ver.Error()})
		return
	}

	_, err := h.svc.CreateDex(ctx.Request.Context(), body)
	if err != nil {
		slog.Error(err.Trace.Error())
		ctx.JSON(err.Code, dto.Response[any]{Message: string(err.Message)})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response[any]{
		Message: "Artifact created successfully",
	})
}
