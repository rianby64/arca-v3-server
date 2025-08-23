package spreadsheet

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getAreasMaterials(ctx context.Context) error {
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

	ranges := "AREAS_MATERIALS!A2:B"
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

	areasMaterialsMap := map[string]*models.AreaMaterials{}
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for index, row := range rowsFromSpreadsheet {
		var (
			material *models.WallMaterial
		)

		areaValue, err := readStringByCellIndex(row, 0)
		if err != nil {
			log.Printf("Skipping row %v: %v", index, err)

			break
		}

		area, err := s.findArea(areaValue)
		if err != nil {
			return errors.Wrapf(err, "error finding area %s in row %v", areaValue, index)
		}

		materialValue := readPtrStringByCellIndex(row, 1)
		if materialValue != nil && *materialValue != "" {
			material, err = s.findMaterial(*materialValue)
			if err != nil {
				return errors.Wrapf(err, "error finding material %s in row %v", *materialValue, index)
			}
		}

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
