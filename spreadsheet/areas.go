package spreadsheet

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"

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

func (s *Spreadsheet) UploadAreasFrom(ctx context.Context, src io.Reader) error {
	var areas models.Areas

	if err := json.NewDecoder(src).Decode(&areas); err != nil {
		return errors.Wrap(err, "Unable to decode areas relations from JSON")
	}

	if len(areas) == 0 {
		return errors.Wrap(models.ErrInvalid, "empty areas relations")
	}

	if err := s.uploadAreas(ctx, areas); err != nil {
		return errors.Wrap(err, "Unable to upload areas relations to spreadsheet")
	}

	return nil
}

func (s *Spreadsheet) uploadAreas(ctx context.Context, areas models.Areas) error {
	requests := []*sheets.Request{}

	for index, area := range areas {
		requests = append(requests, &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          2055988922,
					StartRowIndex:    int64(index) + 1,
					EndRowIndex:      int64(index) + 2,
					StartColumnIndex: 0,
					EndColumnIndex:   3,
				},
				Rows: []*sheets.RowData{
					{
						Values: []*sheets.CellData{
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: &area.Name}},
						},
					},
				},
			},
		})
	}

	if _, err := s.client.Spreadsheets.BatchUpdate(
		s.spreadsheetID,
		&sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).
		Context(ctx).
		Do(); err != nil {
		return err
	}

	return nil
}
