package spreadsheet

import (
	"context"

	"github.com/pkg/errors"

	"arca3/models"
)

func (s *Spreadsheet) getAreasMaterials(ctx context.Context) error {
	if s.materials == nil {
		if err := s.getMaterials(ctx); err != nil {
			return errors.Wrap(err, "Unable to get materials")
		}
	}

	result, err := s.client.Spreadsheets.
		Get(s.spreadsheetID).
		Context(ctx).
		Ranges("AREAS_MATERIALS!A2:B").
		Fields("*").
		IncludeGridData(true).
		Do()
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve spreadsheet")
	}

	areasMaterialsMap := map[string]*models.AreaMaterials{}
	rowsFromSpreadsheet := result.Sheets[0].Data[0].RowData

	for _, row := range rowsFromSpreadsheet {
		var (
			areaValue     *string
			materialValue *string

			material *models.Material
			area     string
		)

		if len(row.Values) > 0 && row.Values[0] != nil && row.Values[0].EffectiveValue != nil && row.Values[0].EffectiveValue.StringValue != nil {
			areaValue = row.Values[0].EffectiveValue.StringValue
		}

		if len(row.Values) > 1 && row.Values[1] != nil && row.Values[1].EffectiveValue != nil && row.Values[1].EffectiveValue.StringValue != nil {
			materialValue = row.Values[1].EffectiveValue.StringValue
		}

		if materialValue != nil {
			for _, m := range s.materials {
				if m.Name == *materialValue {
					material = m
					break
				}
			}
		}

		if areaValue != nil {
			area = *areaValue
		}

		if item, ok := areasMaterialsMap[area]; ok {
			item.Materials = append(item.Materials, material)
		} else {
			areasMaterialsMap[area] = &models.AreaMaterials{
				Area:      area,
				Materials: models.Materials{material},
			}
		}

	}

	areasMaterials := make(models.AreasMaterials, 0, len(areasMaterialsMap))
	for _, item := range areasMaterialsMap {
		areasMaterials = append(areasMaterials, item)
	}

	s.areasMaterials = areasMaterials

	return nil
}
