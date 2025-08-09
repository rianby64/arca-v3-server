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

const (
	spreadsheetID   = "1KAhQuVfNvvsZkEeOcjhrq0zLLQHPq9gNwu5Dpj_VoOY"
	credentialsPath = "arca-v2-account.json"
)

type Spreadsheet struct {
	client *sheets.Service

	materials models.Materials
	// areasMaterials models.AreasMaterials
	// areasKeys      models.AreasKeys
	// relations      models.Relations
}

func New(ctx context.Context) *Spreadsheet {
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	return &Spreadsheet{
		client: client,
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
