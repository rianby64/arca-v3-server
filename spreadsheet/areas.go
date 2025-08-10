package spreadsheet

import (
	"context"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getAreas(ctx context.Context) error {
	ranges := "AREAS!A2:A"
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

	areas := make(models.Areas, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for _, row := range rowsFromSpreadsheet {
		var (
			areaValue *string

			area string
		)

		if len(row.Values) > 0 && row.Values[0] != nil && row.Values[0].EffectiveValue != nil && row.Values[0].EffectiveValue.StringValue != nil {
			areaValue = row.Values[0].EffectiveValue.StringValue
		}

		if areaValue != nil {
			area = *areaValue
		}

		areas = append(areas, &models.Area{
			Name: area,
		})
	}

	s.areas = areas

	return nil
}
