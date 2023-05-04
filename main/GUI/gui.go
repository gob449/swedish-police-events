// Gabriel Räätäri Nyström 2023-05-03
// Hugo Larsson Wilhelmsson 2023-05-03
// This class contains all the GUI-functions that are used in the graphical application
package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	. "project/main/event"
	"strconv"
)

// Initialize and runs the application
func RunGUI() {

	app := app.New()

	openWindow(app)

	app.Run()
}

// Creates and displays the start window
func openWindow(app fyne.App) {
	// start window
	start := app.NewWindow("Start")
	start.CenterOnScreen()
	start.Resize(fyne.NewSize(200, 100))
	contentsOfStartWindow(app, start)

	start.Show()
}

// defines the contents of the start window
func contentsOfStartWindow(app fyne.App, start fyne.Window) {
	//  ok-start button
	okButton := widget.NewButton("ok", func() {
		mainWindow(app)
		start.Close()
	})
	welcomeMessage := widget.NewLabel("Welcome to Swedish Police Events, press ok to get started!")
	startContent := container.New(layout.NewVBoxLayout(),
		welcomeMessage,
		okButton)

	start.SetContent(startContent)
}

// creates and displays the main window
func mainWindow(app fyne.App) {

	// initializing events
	events := AllEventsSlice()

	// creates and displays the main window
	// main window
	windowMain := app.NewWindow("Swedish Police Cases")
	windowMain.Resize(fyne.NewSize(500, 600))
	windowMain.CenterOnScreen()

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fmt.Println("Display settings")
		}),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			fmt.Println("Display search")
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			fmt.Println("Saved events in the archive")
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			fmt.Println("Display help")
		}),
	)

	toolContainer := container.NewVBox(toolbar)

	eventWidget := widget.NewList(
		func() int {
			return len(events)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("label", nil)
		},
		func(i widget.ListItemID, eventButton fyne.CanvasObject) {
			button := eventButton.(*widget.Button)
			button.SetText(events[i].Name)
			button.OnTapped = func() {

				eventInfo := "ID: " + strconv.Itoa(events[i].Id) +
					"\nLocation: " + events[i].Location.Name +
					"\nType: " + events[i].Type +
					"\nSummary: " + events[i].Summary

				eventWindowContainer := container.NewVBox(widget.NewLabel(eventInfo))

				getDescriptionButton := widget.NewButton("Get description (Beware!\nscraping action)", func() {
					eventWindowContainer.Add(widget.NewLabel(GetExtendedSummary(events[i].Url)))
				})

				openInBrowserButton := widget.NewButton("Open in browser", func() {
					OpenInBrowser(events[i].Url)
				})

				closeEventInfoButton := widget.NewButton("Close", func() {
					eventWindowContainer.Hide()
				})

				eventWindowContainer.Add(getDescriptionButton)
				eventWindowContainer.Add(closeEventInfoButton)
				eventWindowContainer.Add(openInBrowserButton)

				toolContainer.Add(eventWindowContainer)
			}
		},
	)

	contentsOfMainWindow := container.NewHSplit(toolContainer, eventWidget)

	windowMain.SetContent(contentsOfMainWindow)

	windowMain.Show()

}
