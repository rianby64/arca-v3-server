package spreadsheet

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getAreasMaterials(ctx context.Context) error {
	if s.materials == nil {
		if err := s.getMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to get materials")
		}
	}

	ranges := "AREAS_MATERIALS!A2:B"
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

	areasMaterialsMap := map[string]*models.AreaMaterials{}
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for index, row := range rowsFromSpreadsheet {
		areaValue, err := readStringByCellIndex(row, 0)
		if err != nil {
			return errors.Wrapf(err, "error reading area name in row %v", index)
		}

		if areaValue == "" {
			return errors.Wrapf(models.ErrInvalid, "empty area name in row %v", index)
		}

		materialValue, err := readStringByCellIndex(row, 1)
		if err != nil {
			return errors.Wrapf(err, "error reading material name in row %v", index)
		}

		if materialValue == "" {
			return errors.Wrapf(models.ErrInvalid, "empty material name in row %v", index)
		}

		material, err := s.findMaterial(materialValue)
		if err != nil {
			return errors.Wrapf(err, "error finding material %s in row %v", materialValue, index)
		}

		area := &models.Area{Name: areaValue}

		if areaMaterial, ok := areasMaterialsMap[area.Name]; ok {
			areaMaterial.Materials = append(areaMaterial.Materials, material)
		} else {
			newAreaMaterial := &models.AreaMaterials{
				Area: area,
			}
			if material != nil {
				newAreaMaterial.Materials = append(newAreaMaterial.Materials, material)
			}
			areasMaterialsMap[area.Name] = newAreaMaterial
		}

	}

	areasMaterials := make(models.AreasMaterials, 0, len(areasMaterialsMap))
	for _, item := range areasMaterialsMap {
		areasMaterials = append(areasMaterials, item)
	}

	s.areasMaterials = areasMaterials

	return nil
}

func (s *Spreadsheet) ReadAreasMaterialsTo(ctx context.Context, dst io.Writer) error {
	if s.areasMaterials == nil {
		if err := s.getAreasMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to read areas materials from spreadsheet")
		}
	}

	if err := json.NewEncoder(dst).Encode(s.areasMaterials); err != nil {
		return errors.Wrap(err, "Unable to encode areas materials to JSON")
	}

	return nil
}
