package handlers

import (
	"context"
	"io"
	"net/http"
)

type Spreadsheet interface {
	ReadAllTo(ctx context.Context, dst io.Writer) error
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
