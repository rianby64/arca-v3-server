package main

import (
	"context"
	"fmt"

	"arca3/config"
	"arca3/spreadsheet"
)

func main() {
	env := config.LoadConfig()
	spreadsheet := spreadsheet.New(context.Background(), env.ServiceCredentialsPath, env.SpreadsheetID)

	// if err := spreadsheet.GetAreasKeys(context.Background()); err != nil {
	// 	fmt.Printf("Error getting areas materials: %v\n", err)
	// 	return
	// }
	fmt.Println(spreadsheet)
}
