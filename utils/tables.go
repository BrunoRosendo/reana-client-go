/*
This file is part of REANA.
Copyright (C) 2022 CERN.

REANA is free software; you can redistribute it and/or modify it
under the terms of the MIT License; see LICENSE file for more details.
*/

package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func DisplayTable(header []interface{}, rows [][]interface{}) {
	rowMatrix := make([]table.Row, len(rows))
	for i, r := range rows {
		rowMatrix[i] = r
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(rowMatrix)
	t.SetStyle(table.StyleRounded)
	t.Render()
}
