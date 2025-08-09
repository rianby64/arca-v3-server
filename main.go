package main

import (
	"context"
	"fmt"

	"arca3/spreadsheet"
)

func main() {
	spreadsheet := spreadsheet.New(context.Background())

	fmt.Println(spreadsheet)
}
