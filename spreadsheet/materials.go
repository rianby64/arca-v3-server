package spreadsheet

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"

	"arca3/models"
)

func (s *Spreadsheet) getMaterials(ctx context.Context) error {
	ranges := "MATERIALS!A2:M"
	result, err := s.client.Spreadsheets.
		Get(s.spreadsheetID).
		Context(ctx).
		Ranges(ranges).
		Fields(effectiveValue).
		IncludeGridData(true).
		Do()
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve spreadsheet %s", ranges)
	}

	materials := make(models.WallMaterials, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for index, row := range rowsFromSpreadsheet {
		material, err := readStringByCellIndex(row, 3)
		if err != nil {
			return errors.Wrapf(err, "error reading material name in row %v", index)
		}

		if material == "" {
			return errors.Wrapf(models.ErrNoData, "empty material name in row %v", index)
		}

		thickness, err := readNumberByCellIndex(row, 1)
		if err != nil {
			return errors.Wrapf(err, "error reading material thickness in row %v", index)
		}

		isStructural, err := readBoolByCellIndex(row, 0)
		if err != nil {
			return errors.Wrapf(err, "error reading isStructural in row %v", index)
		}

		function, err := readStringByCellIndex(row, 2)
		if err != nil {
			return errors.Wrapf(err, "error reading function in row %v", index)
		}

		materials = append(materials, &models.WallMaterial{
			Thickness:    thickness,
			Function:     function,
			IsStructural: isStructural,
			Material: &models.Material{
				Name:                          &material,
				MaterialCategory:              readPtrStringByCellIndex(row, 4),
				CutBackgroundPatternColor:     readPtrStringByCellIndex(row, 5),
				CutBackgroundPatternId:        readPtrStringByCellIndex(row, 6),
				CutForegroundPatternColor:     readPtrStringByCellIndex(row, 7),
				CutForegroundPatternId:        readPtrStringByCellIndex(row, 8),
				SurfaceForegroundPatternColor: readPtrStringByCellIndex(row, 9),
				SurfaceForegroundPatternId:    readPtrStringByCellIndex(row, 10),
				Mark:                          readPtrStringByCellIndex(row, 11),
				Keynote:                       readPtrStringByCellIndex(row, 12),
				Description:                   readPtrStringByCellIndex(row, 13),
				Manufacturer:                  readPtrStringByCellIndex(row, 14),
			},
		})
	}

	s.materials = materials

	return nil
}

func (s *Spreadsheet) ReadMaterialsTo(ctx context.Context, dst io.Writer) error {
	if s.materials == nil {
		if err := s.getMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to read materials from spreadsheet")
		}
	}

	if err := json.NewEncoder(dst).Encode(s.materials); err != nil {
		return errors.Wrap(err, "Unable to encode materials to JSON")
	}

	return nil
}

func (s *Spreadsheet) UploadMaterialsFrom(ctx context.Context, src io.Reader) error {
	var areasRelations models.Materials

	if err := json.NewDecoder(src).Decode(&areasRelations); err != nil {
		return errors.Wrap(err, "Unable to decode areas relations from JSON")
	}

	if len(areasRelations) == 0 {
		return errors.Wrap(models.ErrInvalid, "empty areas relations")
	}

	if err := s.uploadMaterials(ctx, areasRelations); err != nil {
		return errors.Wrap(err, "Unable to upload areas relations to spreadsheet")
	}

	return nil
}

func (s *Spreadsheet) uploadMaterials(ctx context.Context, materials models.Materials) error {
	requests := []*sheets.Request{}

	for index, material := range materials {
		requests = append(requests, &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          1466546092,
					StartRowIndex:    int64(index) + 1,
					EndRowIndex:      int64(index) + 2,
					StartColumnIndex: 0,
					EndColumnIndex:   12,
				},
				Rows: []*sheets.RowData{
					{
						Values: []*sheets.CellData{
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.Name}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.MaterialCategory}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.CutBackgroundPatternColor}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.CutBackgroundPatternId}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.CutForegroundPatternColor}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.CutForegroundPatternId}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.SurfaceForegroundPatternColor}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.SurfaceForegroundPatternId}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.Mark}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.Keynote}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.Description}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: material.Manufacturer}},
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
