package handlers

import (
	"context"
	"io"
	"net/http"
)

type Spreadsheet interface {
	ReadAllTo(ctx context.Context, dst io.Writer) error
	ReadAreasMaterialsTo(ctx context.Context, dst io.Writer) error
	ReadAreasRelationsTo(ctx context.Context, dst io.Writer) error
	ReadAreasTo(ctx context.Context, dst io.Writer) error
	ReadMaterialsTo(ctx context.Context, dst io.Writer) error
}

type WallsHandler struct {
	spreadsheet Spreadsheet
}

func NewWallsHandler(spreadsheet Spreadsheet) *WallsHandler {
	return &WallsHandler{
		spreadsheet: spreadsheet,
	}
}

func (h *WallsHandler) ReadAll(writer http.ResponseWriter, request *http.Request) {
	if err := h.spreadsheet.ReadAllTo(request.Context(), writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *WallsHandler) ReadAreasMaterialsTo(writer http.ResponseWriter, request *http.Request) {
	if err := h.spreadsheet.ReadAreasMaterialsTo(request.Context(), writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *WallsHandler) ReadAreasRelationsTo(writer http.ResponseWriter, request *http.Request) {
	if err := h.spreadsheet.ReadAreasRelationsTo(request.Context(), writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *WallsHandler) ReadAreasTo(writer http.ResponseWriter, request *http.Request) {
	if err := h.spreadsheet.ReadAreasTo(request.Context(), writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *WallsHandler) ReadMaterialsTo(writer http.ResponseWriter, request *http.Request) {
	if err := h.spreadsheet.ReadMaterialsTo(request.Context(), writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
}
