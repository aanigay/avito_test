package segment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateSegment(ctx *gin.Context) {
	req := struct {
		Slug string `json:"slug"`
	}{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "wrong input!")
		return
	}
	segment, err := h.Service.CreateSegment(ctx, req.Slug)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, segment)
}

func (h *Handler) CreateSegmentWithUsers(ctx *gin.Context) {
	var req CreateSegmentWithUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "incorrect input!")
		return
	}
	if req.Percent < 1 || req.Percent > 100 {
		ctx.JSON(http.StatusBadRequest, "incorrect percentage value!")
		return
	}
	err := h.Service.CreateSegmentWithUsers(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) DeleteSegment(ctx *gin.Context) {
	req := struct {
		Slug string `json:"slug"`
	}{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "incorrect input!")
		return
	}
	if err := h.Service.DeleteSegment(ctx, req.Slug); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateSegments(ctx *gin.Context) {
	var req UpdateSegmentsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "incorrect input!")
		return
	}
	if err := h.Service.UpdateSegments(ctx, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) GetSegmentsByUserId(ctx *gin.Context) {
	userId, err := strconv.ParseUint(ctx.Param("userId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "incorrect input!")
		return
	}
	slugs, err := h.Service.GetSegmentsByUserId(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "no such user!")
		return
	}
	ctx.JSON(http.StatusOK, slugs)
}

func (h *Handler) DownloadFile(ctx *gin.Context) {
	ctx.FileAttachment(ctx.Param("filepath")+".csv", ctx.Param("filepath")+".csv")
}

func (h *Handler) GetReports(ctx *gin.Context) {
	match, err := regexp.Match(`^\d{4}-([0]\d|1[0-2])$`, []byte(ctx.Param("date")))
	if !match || err != nil {
		ctx.JSON(http.StatusBadRequest, "wrong date format: please use YYYY-MM")
		return
	}
	url, err := h.Service.GetReports(ctx, ctx.Param("date"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, url)
}
