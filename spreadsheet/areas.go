package spreadsheet

import (
	"context"
	"encoding/json"
	"io"

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

	for index, row := range rowsFromSpreadsheet {
		area, err := readStringByCellIndex(row, 0)
		if err != nil {
			return errors.Wrapf(err, "error reading area name in row %v", index)
		}

		if area == "" {
			return errors.Wrapf(models.ErrInvalid, "empty area name in row %v", index)
		}

		areas = append(areas, &models.Area{
			Name: area,
		})
	}

	s.areas = areas

	return nil
}

func (s *Spreadsheet) ReadAreasTo(ctx context.Context, dst io.Writer) error {
	if s.areas == nil {
		if err := s.getAreas(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas from spreadsheet")
		}
	}

	if err := json.NewEncoder(dst).Encode(s.areas); err != nil {
		return errors.Wrap(err, "Unable to encode areas to JSON")
	}

	return nil
}
