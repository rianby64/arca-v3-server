package spreadsheet

import (
	"context"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) GetAreasKeys(ctx context.Context) error {
	if s.areas == nil {
		if err := s.GetAreas(ctx); err != nil {
			return errors.Wrap(err, "Unable to get areas")
		}
	}

	ranges := "AREAS_KEYS!A2:C"
	result, err := s.client.Spreadsheets.
		Get(s.spreadsheetID).
		Context(ctx).
		Ranges(ranges).
		Fields("*").
		IncludeGridData(true).
		Do()
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve spreadsheet %s", ranges)
	}

	areasKeys := make(models.AreasKeys, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for _, row := range rowsFromSpreadsheet {
		var (
			areaInternalValue *string
			areaExternalValue *string
			keynoteValue      *string

			areaInternal *models.Area
			areaExternal *models.Area
			keynote      string
		)

		if len(row.Values) > 0 && row.Values[0] != nil && row.Values[0].EffectiveValue != nil && row.Values[0].EffectiveValue.StringValue != nil {
			areaInternalValue = row.Values[0].EffectiveValue.StringValue
		}

		if len(row.Values) > 1 && row.Values[1] != nil && row.Values[1].EffectiveValue != nil && row.Values[1].EffectiveValue.StringValue != nil {
			areaExternalValue = row.Values[1].EffectiveValue.StringValue
		}

		if len(row.Values) > 2 && row.Values[2] != nil && row.Values[2].EffectiveValue != nil && row.Values[2].EffectiveValue.StringValue != nil {
			keynoteValue = row.Values[2].EffectiveValue.StringValue
		}

		if areaInternalValue != nil {
			area, err := s.findArea(*areaInternalValue)
			if err != nil {
				return err
			}
			areaInternal = area
		}

		if areaExternalValue != nil {
			area, err := s.findArea(*areaExternalValue)
			if err != nil {
				return err
			}
			areaExternal = area
		}

		if keynoteValue != nil {
			keynote = *keynoteValue
		}

		areasKeys = append(areasKeys, &models.AreaKey{
			AreaInternal: areaInternal,
			AreaExternal: areaExternal,
			Keynote:      keynote,
		})
	}

	s.areasKeys = areasKeys

	return nil
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
