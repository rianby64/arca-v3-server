package spreadsheet

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"arca3/models"
)

type Spreadsheet struct {
	client        *sheets.Service
	spreadsheetID string

	materials models.Materials
	// areasMaterials models.AreasMaterials
	// areasKeys      models.AreasKeys
	// relations      models.Relations
}

func New(ctx context.Context, credentialsPath, spreadsheetID string) *Spreadsheet {
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	return &Spreadsheet{
		client:        client,
		spreadsheetID: spreadsheetID,
	}
}

func (s *Spreadsheet) GetMaterials(ctx context.Context) ([]byte, error) {
	if err := s.getMaterials(ctx); err != nil {
		return nil, err
	}

	data, err := json.Marshal(s.materials)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to marshal materials to JSON")
	}

	return data, nil
}
