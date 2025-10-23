package main

import (
	"datanutsql/drivers"

	"fyne.io/fyne/v2"
)

func ShowMainWindow(a fyne.App, pg *drivers.PG) {
	w := a.NewWindow("DataNutSQL - Query")
	w.Resize(fyne.NewSize(1200, 720))
	//add query funcktions

	w.Show()
}
