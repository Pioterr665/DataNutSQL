package main

import (
	"context"
	"datanutsql/drivers"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowMainWindow(a fyne.App, db drivers.DatabaseDriver) {
	w := a.NewWindow("DataNutSQL - Query")
	w.Resize(fyne.NewSize(1200, 720))

	query_btn := widget.NewButton("Query", nil)
	query_area := widget.NewMultiLineEntry()

	// Get DB Schema przez interfejs
	ctx := context.Background()
	schema, err := db.GetSchemaInfo(ctx)
	if err != nil {
		dialog.ShowError(err, w)
	}

	//Help dialog
	var suggestPopup *widget.PopUp

	sqlKeywords := map[string]struct{}{
		"select": {}, "insert": {}, "update": {}, "delete": {}, "from": {}, "where": {},
		"into": {}, "values": {}, "set": {}, "and": {}, "or": {}, "not": {}, "as": {},
		"on": {}, "join": {}, "left": {}, "right": {}, "inner": {}, "outer": {}, "group": {},
		"by": {}, "order": {}, "limit": {}, "offset": {}, "having": {}, "distinct": {},
		"create": {}, "drop": {}, "alter": {}, "table": {}, "database": {}, "index": {},
		"primary": {}, "key": {}, "foreign": {}, "references": {}, "constraint": {},
		"unique": {}, "check": {}, "default": {}, "null": {}, "is": {}, "true": {}, "false": {},
	}
	query_area.OnChanged = func(text string) {
		if suggestPopup != nil {
			suggestPopup.Hide()
		}
		if len(text) == 0 || schema == nil {
			return
		}
		words := strings.Fields(text)
		if len(words) < 2 {
			return
		}
		prev := strings.ToLower(words[len(words)-2])
		last := strings.ToLower(words[len(words)-1])
		if _, ok := sqlKeywords[prev]; !ok {
			return
		}
		var matches []string
		for table := range schema {
			if strings.HasPrefix(table, last) {
				matches = append(matches, table)
			}
		}
		if len(matches) == 0 {
			return
		}
		list := widget.NewList(
			func() int { return len(matches) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i int, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(matches[i])
			},
		)
		list.OnSelected = func(id int) {
			words[len(words)-1] = matches[id]
			query_area.SetText(strings.Join(words, " "))
			query_area.CursorRow = len(strings.Split(query_area.Text, "\n")) - 1
			query_area.CursorColumn = len(words[len(words)-1])
			if suggestPopup != nil {
				suggestPopup.Hide()
			}
		}

		// Oblicz pozycję popupu względem query_area
		cursorRow := len(strings.Split(query_area.Text[:query_area.CursorColumn], "\n")) - 1
		cursorCol := query_area.CursorColumn
		offsetY := float32(24 + cursorRow*18)
		offsetX := float32(8 + cursorCol*8)

		suggestPopup = widget.NewPopUp(list, w.Canvas())
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(query_area)
		suggestPopup.Move(fyne.NewPos(pos.X+offsetX, pos.Y+offsetY))
		suggestPopup.Resize(fyne.NewSize(200, 200))
		suggestPopup.Show()
	}

	var tableNames []string
	for table := range schema {
		tableNames = append(tableNames, table)
	}

	tree := widget.NewTree(
		func(uid string) []string {
			if uid == "" {
				return tableNames
			}
			return nil
		},
		func(uid string) bool {
			return uid == ""
		},
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(uid string, branch bool, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(uid)
		},
	)
	tree.OpenAllBranches()

	var tableData [][]any
	table := widget.NewTable(
		func() (int, int) {
			if len(tableData) == 0 {
				return 1, 1
			}
			return len(tableData), len(tableData[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			i, j := id.Row, id.Col
			if len(tableData) == 0 {
				label.SetText("")
			} else {
				label.SetText(fmt.Sprintf("%v", tableData[i][j]))
			}
		},
	)
	minWidth := float32(120)

	query_btn.OnTapped = func() {
		sql := query_area.SelectedText()
		if sql == "" {
			sql = query_area.Text
		}
		if strings.TrimSpace(sql) == "" {
			return
		}
		ctx := context.Background()
		result, err := db.QueryTableWithHeaders(ctx, sql)
		if err != nil {
			tableData = [][]any{{"Error:", err.Error()}}
		} else if len(result) == 0 {
			tableData = [][]any{{"No results"}}
		} else {
			tableData = result
		}
		if len(tableData) > 0 {
			for col := 0; col < len(tableData[0]); col++ {
				table.SetColumnWidth(col, minWidth)
			}
		}
		table.Refresh()
	}

	top := container.NewBorder(
		query_btn,
		nil, nil, nil,
		query_area,
	)

	tableScroll := container.NewVScroll(table)

	split := container.NewVSplit(top, tableScroll)
	split.Offset = 0.25
	treeInfo := widget.NewLabel("Available tables:")
	leftPanel := container.NewBorder(treeInfo, nil, nil, nil, tree)
	content := container.NewHSplit(
		container.NewVScroll(leftPanel),
		split,
	)
	content.Offset = 0.2
	w.SetContent(content)
	w.Show()
}
