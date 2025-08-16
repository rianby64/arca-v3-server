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

	materials      models.WallMaterials
	areas          models.Areas
	areasMaterials models.AreasMaterials
	relations      models.AreasRelations
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

func readStringByCellIndex(row *sheets.RowData, index int) (string, error) {
	if len(row.Values) <= index {
		return "", errors.Wrapf(models.ErrInvalid, "index %d out of range for row with %d values", index, len(row.Values))
	}

	if row.Values[index] == nil {
		return "", errors.Wrapf(models.ErrInvalid, "no value at index %d in row", index)
	}

	if row.Values[index].EffectiveValue == nil {
		return "", errors.Wrapf(models.ErrNoData, "no effective value at index %d in row", index)
	}

	if row.Values[index].EffectiveValue.StringValue == nil {
		return "", errors.Wrapf(models.ErrInvalid, "no string value at index %d in row", index)
	}

	value := row.Values[index].EffectiveValue.StringValue

	return *value, nil
}

func readNumberByCellIndex(row *sheets.RowData, index int) (float64, error) {
	if len(row.Values) <= index {
		return 0, errors.Wrapf(models.ErrInvalid, "index %d out of range for row with %d values", index, len(row.Values))
	}

	if row.Values[index] == nil {
		return 0, errors.Wrapf(models.ErrInvalid, "no value at index %d in row", index)
	}

	if row.Values[index].EffectiveValue == nil {
		return 0, errors.Wrapf(models.ErrNoData, "no effective value at index %d in row", index)
	}

	if row.Values[index].EffectiveValue.NumberValue == nil {
		return 0, errors.Wrapf(models.ErrInvalid, "no number value at index %d in row", index)
	}

	value := row.Values[index].EffectiveValue.NumberValue

	return *value, nil
}

func readBoolByCellIndex(row *sheets.RowData, index int) (bool, error) {
	if len(row.Values) <= index {
		return false, errors.Wrapf(models.ErrInvalid, "index %d out of range for row with %d values", index, len(row.Values))
	}

	if row.Values[index] == nil {
		return false, errors.Wrapf(models.ErrInvalid, "no value at index %d in row", index)
	}

	if row.Values[index].EffectiveValue == nil {
		return false, errors.Wrapf(models.ErrNoData, "no effective value at index %d in row", index)
	}

	if row.Values[index].EffectiveValue.BoolValue == nil {
		return false, errors.Wrapf(models.ErrInvalid, "no boolean value at index %d in row", index)
	}

	value := row.Values[index].EffectiveValue.BoolValue

	return *value, nil
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

func (s *Spreadsheet) findMaterial(name string) (*models.WallMaterial, error) {
	if name == "" {
		return nil, errors.Wrapf(models.ErrInvalid, "empty material name")
	}

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
	if s.relations == nil {
		if err := s.getAreasRelations(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas relations from spreadsheet")
		}
	}

	allEntries := map[string]any{
		"Materials":      s.materials,
		"Areas":          s.areas,
		"AreasMaterials": s.areasMaterials,
		"Relations":      s.relations,
	}

	if err := json.NewEncoder(dst).Encode(allEntries); err != nil {
		return errors.Wrap(err, "Unable to encode all entries to JSON")
	}

	return nil
}
