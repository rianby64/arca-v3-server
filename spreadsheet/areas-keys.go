package spreadsheet

import (
	"context"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getAreasKeys(ctx context.Context) error {
	result, err := s.client.Spreadsheets.
		Get(s.spreadsheetID).
		Context(ctx).
		Ranges("AREAS_KEYS!A2:C").
		Fields("*").
		IncludeGridData(true).
		Do()
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve spreadsheet")
	}

	areasKeys := make(models.AreasKeys, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for _, row := range rowsFromSpreadsheet {
		var (
			areaInternalValue *string
			areaExternalValue *string
			keynoteValue      *string

			areaInternal string
			areaExternal string
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
			areaInternal = *areaInternalValue
		}

		if areaExternalValue != nil {
			areaExternal = *areaExternalValue
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
