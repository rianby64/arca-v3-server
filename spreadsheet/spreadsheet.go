package spreadsheet

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"arca3/models"
)

type Spreadsheet struct {
	client        *sheets.Service
	spreadsheetID string

	materials      models.Materials
	areas          models.Areas
	areasMaterials models.AreasMaterials
	areasKeys      models.AreasKeys
	relations      models.Relations
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

func (s *Spreadsheet) findArea(name string) (*models.Area, error) {
	if s.areas == nil {
		return nil, models.ErrUnavailable
	}

	for _, area := range s.areas {
		if area.Name == name {
			return area, nil
		}
	}

	return nil, errors.Wrapf(models.ErrNotFound, "area %s", name)
}

func (s *Spreadsheet) findMaterial(name string) (*models.Material, error) {
	if s.materials == nil {
		return nil, models.ErrUnavailable
	}

	for _, material := range s.materials {
		if material.Name == name {
			return material, nil
		}
	}

	return nil, errors.Wrapf(models.ErrNotFound, "material %s", name)
}

func (s *Spreadsheet) ReadAllTo(ctx context.Context, dst io.Writer) error {
	if s.materials == nil {
		if err := s.getMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to read materials from spreadsheet")
		}
	}
	if s.areas == nil {
		if err := s.getAreas(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas from spreadsheet")
		}
	}
	if s.areasMaterials == nil {
		if err := s.getAreasMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas materials from spreadsheet")
		}
	}
	if s.areasKeys == nil {
		if err := s.getAreasKeys(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas keys from spreadsheet")
		}
	}
	if s.relations == nil {
		if err := s.getAreasRelations(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas relations from spreadsheet")
		}
	}

	allEntries := map[string]any{
		"materials":       s.materials,
		"areas":           s.areas,
		"areas_materials": s.areasMaterials,
		"areas_keys":      s.areasKeys,
		"relations":       s.relations,
	}

	if err := json.NewEncoder(dst).Encode(allEntries); err != nil {
		return errors.Wrap(err, "Unable to encode all entries to JSON")
	}

	return nil
}
