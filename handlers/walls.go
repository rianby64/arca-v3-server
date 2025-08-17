package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
)

type Spreadsheet interface {
	ResetData()

	ReadAreasMaterialsTo(ctx context.Context, dst io.Writer) error

	ReadAreasRelationsTo(ctx context.Context, dst io.Writer) error
	UploadAreasRelationsFrom(ctx context.Context, dst io.Reader) error

	ReadAreasTo(ctx context.Context, dst io.Writer) error
	UploadAreasFrom(ctx context.Context, dst io.Reader) error

	ReadMaterialsTo(ctx context.Context, dst io.Writer) error
	UploadMaterialsFrom(ctx context.Context, src io.Reader) error
}

type WallsHandler struct {
	spreadsheet Spreadsheet
}

func NewWallsHandler(spreadsheet Spreadsheet) *WallsHandler {
	return &WallsHandler{
		spreadsheet: spreadsheet,
	}
}

func (h *WallsHandler) ResetData(writer http.ResponseWriter, request *http.Request) {
	h.spreadsheet.ResetData()
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

func (h *WallsHandler) UploadAreasRelationsFrom(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if err := h.spreadsheet.UploadAreasRelationsFrom(request.Context(), request.Body); err != nil {
		log.Printf("Error uploading areas relations: %v", err)
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

func (h *WallsHandler) UploadAreasFrom(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if err := h.spreadsheet.UploadAreasFrom(request.Context(), request.Body); err != nil {
		log.Printf("Error uploading areas: %v", err)
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

func (h *WallsHandler) UploadMaterialsFrom(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if err := h.spreadsheet.UploadMaterialsFrom(request.Context(), request.Body); err != nil {
		log.Printf("Error uploading materials: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
}
