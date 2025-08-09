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

	fmt.Println(spreadsheet)
}
