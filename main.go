package main

import (
	"datanutsql/drivers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("DataNutSQL")
	w.Resize(fyne.NewSize(600, 400))

	//login inputs
	address := widget.NewEntry()
	address.SetPlaceHolder("DB Address")

	port := widget.NewEntry()
	port.SetPlaceHolder("DB port")

	db := widget.NewEntry()
	db.SetPlaceHolder("Db Name")

	user := widget.NewEntry()
	user.SetPlaceHolder("Username")

	password := widget.NewEntry()
	password.SetPlaceHolder("Password")

	login_btn := widget.NewButton("LogIn", nil)

	status := widget.NewLabel("")

	login_btn.OnTapped = func() {
		status.SetText("Connecting...")
		//conn_string := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user.Text, password.Text, address.Text, port.Text, db.Text)
		//for testing:
		conn_string := "postgres://postgres:postgres@localhost:5432/postgres"
		pg, err := drivers.ConnectPG(conn_string)
		if err != nil {
			status.SetText("Error")
		} else {

			status.SetText("Connected!")
			ShowMainWindow(a, pg)
			w.Close()

		}

	}

	content := container.NewVBox(address,
		port,
		db,
		user,
		password,
		login_btn,
		status,
	)
	w.SetContent(content)
	w.ShowAndRun()
}
