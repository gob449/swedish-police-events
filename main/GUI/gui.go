// Gabriel Räätäri Nyström 2023-05-03
// Hugo Larsson Wilhelmsson 2023-05-03
// This class contains all the GUI-functions that are used in the graphical application
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fmt"
	"log"
)

func RunGUI() {
	app := app.New()

	openWindow(app)

	app.Run()
}

func openWindow(app fyne.App) {
	// start window
	start := app.NewWindow("Start")
	start.CenterOnScreen()
	start.Resize(fyne.NewSize(200, 100))
	contentsOfStartWindow(app, start)

	start.Show()
}

func contentsOfStartWindow(app fyne.App, start fyne.Window) {
	//  ok-start button
	okButton := widget.NewButton("ok", func() {
		mainWindow(app)
		start.Close()
	})

	myTextStyle := fyne.TextStyle{Bold: true, Monospace: true}
	welcomeMessage := widget.NewLabel("Welcome to Swedish Police Events, press ok to get started!")
	welcomeMessage.TextStyle = myTextStyle

	startContent := container.New(layout.NewVBoxLayout(),
		welcomeMessage,
		okButton)

	start.SetContent(startContent)
}

func mainWindow(app fyne.App) {
	// main window
	windowMain := app.NewWindow("SwePe")
	windowMain.Resize(fyne.NewSize(500, 600))
	windowMain.CenterOnScreen()

	contentsOfMainWindow := container.New(layout.NewVBoxLayout(), verticalToolbar())

	windowMain.SetContent(contentsOfMainWindow)

	windowMain.Show()

}

func verticalToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fmt.Println("Display settings")
		}),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			fmt.Println("Display search")
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			fmt.Println("Display save")
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	return toolbar
}