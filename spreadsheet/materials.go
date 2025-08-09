package spreadsheet

import (
	"context"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getMaterials(ctx context.Context) error {
	result, err := s.client.Spreadsheets.
		Get(s.spreadsheetID).
		Context(ctx).
		Ranges("MATERIALS!A2:O").
		Fields("*").
		IncludeGridData(true).
		Do()
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve spreadsheet")
	}

	materials := make(models.Materials, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for _, row := range rowsFromSpreadsheet {
		var (
			materialValue     *string
			keynoteValue      *string
			thicknessValue    *float64
			isStructuralValue *bool

			material     string
			keynote      string
			thickness    float64
			isStructural bool
		)

		if len(row.Values) > 0 && row.Values[0] != nil && row.Values[0].EffectiveValue != nil && row.Values[0].EffectiveValue.StringValue != nil {
			materialValue = row.Values[0].EffectiveValue.StringValue
		}

		if len(row.Values) > 1 && row.Values[1] != nil && row.Values[1].EffectiveValue != nil && row.Values[1].EffectiveValue.NumberValue != nil {
			thicknessValue = row.Values[1].EffectiveValue.NumberValue
		}

		if len(row.Values) > 2 && row.Values[2] != nil && row.Values[2].EffectiveValue != nil && row.Values[2].EffectiveValue.BoolValue != nil {
			isStructuralValue = row.Values[2].EffectiveValue.BoolValue
		}

		if len(row.Values) > 3 && row.Values[3] != nil && row.Values[3].EffectiveValue != nil && row.Values[3].EffectiveValue.StringValue != nil {
			keynoteValue = row.Values[3].EffectiveValue.StringValue
		}

		if materialValue != nil {
			material = *materialValue
		} else {
			material = ""
		}

		if keynoteValue != nil {
			keynote = *keynoteValue
		} else {
			keynote = ""
		}

		if thicknessValue != nil {
			thickness = *thicknessValue
		} else {
			thickness = 0.0
		}

		if isStructuralValue != nil {
			isStructural = *isStructuralValue
		} else {
			isStructural = false
		}

		materials = append(materials, &models.Material{
			Name:         material,
			Thickness:    thickness,
			Keynote:      keynote,
			IsStructural: isStructural,
		})
	}

	s.materials = materials

	return nil
}
