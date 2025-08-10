package spreadsheet

import (
	"context"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getAreasRelations(ctx context.Context) error {
	if s.areas == nil {
		if err := s.getAreas(ctx); err != nil {
			return errors.Wrap(err, "Unable to get areas")
		}
	}
	if s.materials == nil {
		if err := s.getMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to get materials")
		}
	}

	ranges := "AREAS_RELATIONS!A2:D"
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

	areasKeys := make(models.Relations, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for _, row := range rowsFromSpreadsheet {
		var (
			areaInternalValue *string
			areaExternalValue *string
			materialValue     *string
			sameAreaValue     *bool

			areaInternal *models.Area
			areaExternal *models.Area
			material     *models.Material
			sameArea     bool
		)

		if len(row.Values) > 0 && row.Values[0] != nil && row.Values[0].EffectiveValue != nil && row.Values[0].EffectiveValue.StringValue != nil {
			areaInternalValue = row.Values[0].EffectiveValue.StringValue
		}

		if len(row.Values) > 1 && row.Values[1] != nil && row.Values[1].EffectiveValue != nil && row.Values[1].EffectiveValue.StringValue != nil {
			areaExternalValue = row.Values[1].EffectiveValue.StringValue
		}

		if len(row.Values) > 2 && row.Values[2] != nil && row.Values[2].EffectiveValue != nil && row.Values[2].EffectiveValue.StringValue != nil {
			materialValue = row.Values[2].EffectiveValue.StringValue
		}

		if len(row.Values) > 3 && row.Values[3] != nil && row.Values[3].EffectiveValue != nil && row.Values[3].EffectiveValue.BoolValue != nil {
			sameAreaValue = row.Values[3].EffectiveValue.BoolValue
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

		if materialValue != nil {
			materialFound, err := s.findMaterial(*materialValue)
			if err != nil {
				return err
			}
			material = materialFound
		}

		if sameAreaValue != nil {
			sameArea = *sameAreaValue
		}

		areasKeys = append(areasKeys, &models.Relation{
			AreaInternal: areaInternal,
			AreaExternal: areaExternal,
			Material:     material,
			SameArea:     sameArea,
		})
	}

	s.relations = areasKeys

	return nil
}
