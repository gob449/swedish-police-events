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
	. "project/main/event"
	"strconv"
	"time"
)

// Initialize and runs the application
func RunGUI() {

	SwePe := app.New()

	openWindow(SwePe)

	SwePe.Run()
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

func mainWindow(app fyne.App) {

	mainWindow := app.NewWindow("Swedish Police Events")
	mainWindow.Resize(fyne.NewSize(1200, 700))
	mainWindow.CenterOnScreen()

	allEvents := AllEventsSlice()
	allEventsList := eventListView(allEvents)

	eventInfo := widget.NewLabel("Please select an event")
	eventInfo.Wrapping = fyne.TextWrapWord
	extensiveSummary := widget.NewLabel("")
	extensiveSummary.Wrapping = fyne.TextWrapWord
	openInBrowserButton := widget.NewButton("Open in browser", func() {})
	openInBrowserButton.Hide()
	scrapeBrowserButton := widget.NewButton("Scrape webpage for summary", func() {})
	scrapeBrowserButton.Hide()

	displayEventInfo := container.NewVBox(eventInfo, extensiveSummary, openInBrowserButton, scrapeBrowserButton)
	allEventsList.OnSelected = eventOnSelection(allEvents, eventInfo, extensiveSummary, openInBrowserButton, scrapeBrowserButton)

	typeMenuOptions := keysListview(TypeKeys)

	// THIS METHOD CAN BE REUSED FOR LOCATIONMENU OPTIOS FOR INSTANCE.
	// BY PROVIDING SEARCHKEYS
	typeMenuOptions.OnSelected = func(id widget.ListItemID) {
		typeKey := TypeKeys[id]
		// Creating and editing the window that pops up when a type has been chosen
		subCatWindow := app.NewWindow(typeKey)
		subCatWindow.Resize(fyne.NewSize(400, 400))
		subCatWindow.CenterOnScreen()
		subCatEvents := SubCatType(AllEventsSlice(), typeKey)
		subCatEventsListView := eventListView(subCatEvents)
		subCatEventsListView.OnSelected = eventOnSelection(subCatEvents, eventInfo, extensiveSummary, openInBrowserButton, scrapeBrowserButton)
		subCatWindow.SetContent(subCatEventsListView)
		subCatWindow.Show()
	}

	typeMenuOptionsPopUp := widget.NewPopUp(typeMenuOptions, mainWindow.Canvas())
	typeMenuOptionsPopUp.Resize(fyne.NewSize(300, 200))
	typeSearch := fyne.NewMenuItem("Type", func() {
		typeMenuOptionsPopUp.Show()
	})

	locationMenuOptions := keysListview(LocationKeys)
	locationMenuOptions.OnSelected = func(id widget.ListItemID) {
		locationKey := LocationKeys[id]
		// Creating and editing the window that pops up when a type has been chosen
		subCatWindow := app.NewWindow(locationKey)
		subCatWindow.Resize(fyne.NewSize(400, 400))
		subCatWindow.CenterOnScreen()
		subCatEvents := SubCatLocation(AllEventsSlice(), locationKey)
		subCatEventsListView := eventListView(subCatEvents)
		subCatEventsListView.OnSelected = eventOnSelection(subCatEvents, eventInfo, extensiveSummary, openInBrowserButton, scrapeBrowserButton)
		subCatWindow.SetContent(subCatEventsListView)
		subCatWindow.Show()
	}
	locationMenuOptionsPopUp := widget.NewPopUp(locationMenuOptions, mainWindow.Canvas())
	locationMenuOptionsPopUp.Resize(fyne.NewSize(300, 200))
	locationSearch := fyne.NewMenuItem("Location", func() {
		locationMenuOptionsPopUp.Show()
	})

	searchMenu := fyne.NewMenu("Search", typeSearch, locationSearch)
	searchMenuPopUp := widget.NewPopUpMenu(searchMenu, mainWindow.Canvas())

	saveButton := widget.NewButton("Save", func() {
		SaveInArchive(AllEventsSlice())

		notificationMessage := widget.NewLabel("Updating archive...")
		saveNotification := widget.NewPopUp(notificationMessage, mainWindow.Canvas())
		saveNotification.Show()
		select {
		case <-time.After(2 * time.Second):
		}
		notificationMessage.SetText("Archive has been updated")
		select {
		case <-time.After(2 * time.Second):
		}
		saveNotification.Hide()
	})

	saveMessage := widget.NewLabel("Click 'Save' to update current archive")
	savePopUpContainer := container.NewVBox(saveMessage, saveButton)
	savePopUp := widget.NewPopUp(savePopUpContainer, mainWindow.Canvas())

	verticalToolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			searchMenuPopUp.Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			savePopUp.Resize(fyne.NewSize(200, 100))
			savePopUp.Show()
		}),
	)

	eventsListAndInfoDisplay := container.NewHSplit(allEventsList, container.NewMax(displayEventInfo))
	mainWindowContainer := container.NewVSplit(verticalToolbar, eventsListAndInfoDisplay)
	mainWindowContainer.SetOffset(0.05)

	mainWindow.SetContent(mainWindowContainer)
	mainWindow.Show()
}

/*
	By providing a slice of valid searchKeys the method provides a listview of the keys

where all alternatives are visible.
*/
func keysListview(keys []string) *widget.List {
	return widget.NewList(
		func() int {
			return len(keys)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Type options")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			event := object.(*widget.Label)
			event.SetText(keys[id])
		},
	)
}

/*
	The consequences of an event being selected, the parameter events []event is all events

or also a subcategory of events, this method helps to bring a standard layout
to the program
*/
func eventOnSelection(events []Event, eventInfo *widget.Label, extensiveSummary *widget.Label, openInBrowserButton *widget.Button, scrapeBrowserButton *widget.Button) func(id widget.ListItemID) {
	return func(id widget.ListItemID) {
		info := "ID: " + strconv.Itoa(events[id].Id) +
			"\nLocation: " + events[id].Location.Name +
			"\nType: " + events[id].Type +
			"\nSummary: " + events[id].Summary
		eventInfo.SetText(info)
		extensiveSummary.SetText("")
		openInBrowserButton.Show()
		openInBrowserButton.OnTapped = func() {
			OpenInBrowser(events[id].Url)
		}
		scrapeBrowserButton.Show()
		scrapeBrowserButton.OnTapped = func() {
			extensiveSummary.SetText(GetExtendedSummary(events[id].Url))
		}
	}
}

func eventListView(events []Event) *widget.List {
	eventsList := widget.NewList(
		func() int {
			return len(events)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Event")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			event := object.(*widget.Label)
			event.SetText(events[id].Name)
		},
	)
	return eventsList
}
