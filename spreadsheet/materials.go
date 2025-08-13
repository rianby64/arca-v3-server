package spreadsheet

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getMaterials(ctx context.Context) error {
	ranges := "MATERIALS!A2:M"
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

	materials := make(models.Materials, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for index, row := range rowsFromSpreadsheet {
		material, err := readStringByCellIndex(row, 1)
		if err != nil {
			return errors.Wrapf(err, "error reading material name in row %v", index)
		}

		if material == "" {
			return errors.Wrapf(models.ErrNoData, "empty material name in row %v", index)
		}

		thickness, err := readNumberByCellIndex(row, 2)
		if err != nil {
			return errors.Wrapf(err, "error reading material thickness in row %v", index)
		}

		isStructural, err := readBoolByCellIndex(row, 0)
		if err != nil {
			return errors.Wrapf(err, "error reading isStructural in row %v", index)
		}

		keynote, err := readStringByCellIndex(row, 9)
		if err != nil {
			return errors.Wrapf(err, "error reading keynote in row %v", index)
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
