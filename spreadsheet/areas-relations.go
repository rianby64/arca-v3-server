package spreadsheet

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"

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

	ranges := "AREAS_RELATIONS!A2:E"
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

	areasKeys := make(models.AreasRelations, 0, len(result.Sheets[0].Data[0].RowData))
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for index, row := range rowsFromSpreadsheet {
		areaInternalValue, err := readStringByCellIndex(row, 2)
		if err != nil {
			return errors.Wrapf(err, "error reading area internal name in row %v", index)
		}

		if areaInternalValue == "" {
			return errors.Wrapf(models.ErrInvalid, "empty area internal name in row %v", index)
		}

		areaInternal, err := s.findArea(areaInternalValue)
		if err != nil {
			return errors.Wrapf(err, "error finding area %s in row %v", areaInternalValue, index)
		}

		areaExternalValue, err := readStringByCellIndex(row, 3)
		if err != nil {
			return errors.Wrapf(err, "error reading area external name in row %v", index)
		}

		if areaExternalValue == "" {
			return errors.Wrapf(models.ErrInvalid, "empty area external name in row %v", index)
		}

		areaExternal, err := s.findArea(areaExternalValue)
		if err != nil {
			return errors.Wrapf(err, "error finding area %s in row %v", areaExternalValue, index)
		}

		materialValue, _ := readStringByCellIndex(row, 4)

		material, err := s.findMaterial(materialValue)
		if !errors.Is(err, models.ErrInvalid) && err != nil {
			return errors.Wrapf(err, "error finding material %s in row %v", materialValue, index)
		}

		sameArea, err := readBoolByCellIndex(row, 0)
		if err != nil {
			return errors.Wrapf(err, "error reading sameArea in row %v", index)
		}

		wallKeynote, err := readStringByCellIndex(row, 1)
		if !errors.Is(err, models.ErrNoData) && err != nil {
			return errors.Wrapf(err, "error reading wallKeynote in row %v", index)
		}

		areasKeys = append(areasKeys, &models.AreaRelation{
			AreaInternal: areaInternal,
			AreaExternal: areaExternal,
			Material:     material,
			SameArea:     sameArea,
			WallKeynote:  wallKeynote,
		})
	}

	s.relations = areasKeys

	return nil
}

func (s *Spreadsheet) ReadAreasRelationsTo(ctx context.Context, dst io.Writer) error {
	if s.relations == nil {
		if err := s.getAreasRelations(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas relations from spreadsheet")
		}
	}

	if err := json.NewEncoder(dst).Encode(s.relations); err != nil {
		return errors.Wrap(err, "Unable to encode areas relations to JSON")
	}

	return nil
}

func (s *Spreadsheet) UploadAreasRelationsFrom(ctx context.Context, src io.Reader) error {
	var areasRelations models.AreasRelations

	if err := json.NewDecoder(src).Decode(&areasRelations); err != nil {
		return errors.Wrap(err, "Unable to decode areas relations from JSON")
	}

	if len(areasRelations) == 0 {
		return errors.Wrap(models.ErrInvalid, "empty areas relations")
	}

	if err := s.uploadAreasRelations(ctx, areasRelations); err != nil {
		return errors.Wrap(err, "Unable to upload areas relations to spreadsheet")
	}

	return nil
}

func (s *Spreadsheet) uploadAreasRelations(ctx context.Context, areasRelations models.AreasRelations) error {
	requests := []*sheets.Request{}

	for index, relation := range areasRelations {
		var (
			areaInternal, areaExternal *string
		)

		if relation.AreaInternal != nil {
			areaInternal = &relation.AreaInternal.Name
		}

		if relation.AreaExternal != nil {
			areaExternal = &relation.AreaExternal.Name
		}

		requests = append(requests, &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          1715124245,
					StartRowIndex:    int64(index) + 1,
					EndRowIndex:      int64(index) + 2,
					StartColumnIndex: 0,
					EndColumnIndex:   3,
				},
				Rows: []*sheets.RowData{
					{
						Values: []*sheets.CellData{
							{UserEnteredValue: &sheets.ExtendedValue{BoolValue: &relation.SameArea}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: areaInternal}},
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: areaExternal}},
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
