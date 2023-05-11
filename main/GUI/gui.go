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
	"strings"
	"time"
)

// Initialize and runs the application
func RunGUI() {

	SwePe := app.New()

	// Opens start window
	openWindow(SwePe)

	SwePe.Run()
}

// Creates and displays the start window
func openWindow(app fyne.App) {
	// start window
	start := app.NewWindow("Start")
	start.CenterOnScreen() // center start window
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

// mainWindow creates and defines a Fyne 'fyne.Window' object
func mainWindow(app fyne.App) {

	mainWindow := app.NewWindow("Swedish Police Events")
	mainWindow.Resize(fyne.NewSize(1200, 700))
	mainWindow.CenterOnScreen()

	allEvents := AllEventsSlice()
	allEventsList := eventListView(allEvents)

	/*
		The layout that is displayed when an ecent is selected
	*/
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

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/*
		Segment holds script for the submenu "Type" under "Search" toolbar option
	*/

	// Keeps track of keys matching query
	filteredTypeOptions := make([]string, 0)

	// Create a typeMenuOptions widget to display the items
	typeMenuOptions := keysListview(TypeKeys)

	// Entry widget for search query
	typeSearchEntry := widget.NewEntry()
	typeSearchEntry.SetPlaceHolder("Search")

	// Event handler for search query entry
	typeSearchEntry.OnChanged = func(query string) {
		filteredTypeElements := make([]string, 0)

		// If search query is empty, show all items
		if query == "" {
			filteredTypeElements = TypeKeys
		} else {
			// Filter items based on search query
			for _, item := range TypeKeys {
				if containsIgnoreCase(item, query) {
					filteredTypeElements = append(filteredTypeElements, item)
				}
			}
		}

		// Update the list with the filtered items
		typeMenuOptions.Length = func() int {
			return len(filteredTypeElements)
		}
		typeMenuOptions.CreateItem = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}
		typeMenuOptions.UpdateItem = func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(filteredTypeElements[i])
		}
		typeMenuOptions.Refresh()
		filteredTypeOptions = filteredTypeElements
	}

	// THIS METHOD CAN BE REUSED FOR LOCATIONMENU OPTIOS FOR INSTANCE.
	// BY PROVIDING SEARCHKEYS
	typeMenuOptions.OnSelected = func(id widget.ListItemID) {
		var typeKey string
		if len(filteredTypeOptions) == 0 {
			typeKey = TypeKeys[id]
		} else {
			typeKey = filteredTypeOptions[id]
		}
		// Creating and editing the window that pops up when a type has been chosen
		subCatWindow := app.NewWindow(typeKey)
		subCatWindow.Resize(fyne.NewSize(400, 400))
		subCatWindow.CenterOnScreen()
		subCatEvents := SubCatType(allEvents, typeKey) // Creates a []Events of the type subcategory matching the key
		subCatEventsListView := eventListView(subCatEvents)
		subCatEventsListView.OnSelected = eventOnSelection(subCatEvents, eventInfo, extensiveSummary, openInBrowserButton, scrapeBrowserButton)
		subCatWindow.SetContent(subCatEventsListView)
		subCatWindow.Show()
	}

	typeSearchEntryAndOptions := container.NewVSplit(typeSearchEntry, typeMenuOptions)
	typeSearchEntryAndOptions.SetOffset(0.1)

	typeMenuOptionsPopUp := widget.NewPopUp(typeSearchEntryAndOptions, mainWindow.Canvas())
	typeMenuOptionsPopUp.Resize(fyne.NewSize(250, 250))
	typeSearch := fyne.NewMenuItem("Type", func() {
		typeMenuOptionsPopUp.Show()
	})

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	/*
		Segment holds script for the submenu "Location" under "Search" toolbar option, identical to type search menu see earlier segemt
	*/

	locationMenuOptions := keysListview(LocationKeys)
	filteredLocationOptions := make([]string, 0)

	// Entry widget for search query
	locationSearchEntry := widget.NewEntry()
	locationSearchEntry.SetPlaceHolder("Enter Location")

	// Event handler for search query entry
	locationSearchEntry.OnChanged = func(query string) {
		filteredLocationElements := make([]string, 0)

		// If search query is empty, show all items
		if query == "" {
			filteredLocationElements = LocationKeys
		} else {
			// Filter items based on search query
			for _, item := range LocationKeys {
				if containsIgnoreCase(item, query) {
					filteredLocationElements = append(filteredLocationElements, item)
				}
			}
		}

		// Update the list with the filtered items
		locationMenuOptions.Length = func() int {
			return len(filteredLocationElements)
		}
		locationMenuOptions.CreateItem = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}
		locationMenuOptions.UpdateItem = func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(filteredLocationElements[i])
		}
		locationMenuOptions.Refresh()
		filteredLocationOptions = filteredLocationElements
	}

	locationMenuOptions.OnSelected = func(id widget.ListItemID) {
		var locationKey string
		if len(filteredLocationOptions) == 0 {
			locationKey = LocationKeys[id]
		} else {
			locationKey = filteredLocationOptions[id]
		}
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

	locSearchEntryAndOptions := container.NewVSplit(locationSearchEntry, locationMenuOptions)
	locSearchEntryAndOptions.SetOffset(0.1)

	locSearchEntryAndOptionsPopUp := widget.NewPopUp(locSearchEntryAndOptions, mainWindow.Canvas())
	locSearchEntryAndOptionsPopUp.Resize(fyne.NewSize(250, 250))
	locationSearch := fyne.NewMenuItem("Location", func() {
		locSearchEntryAndOptionsPopUp.Show()
	})

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

	darkThemeOption := fyne.NewMenuItem("Dark theme", func() {
		darkTheme := theme.DarkTheme()
		app.Settings().SetTheme(darkTheme)
	})

	lightThemeOption := fyne.NewMenuItem("Light theme", func() {
		lightTheme := theme.LightTheme()
		app.Settings().SetTheme(lightTheme)
	})

	themeOptionsMenu := fyne.NewMenu("Themes", darkThemeOption, lightThemeOption)
	themeOptionsMenuPopUp := widget.NewPopUpMenu(themeOptionsMenu, mainWindow.Canvas())

	themeOption := fyne.NewMenuItem("Theme Option", func() {
		themeOptionsMenuPopUp.Show()
	})

	plusButton := widget.NewButton("+", func() {
		currentSize := mainWindow.Canvas().Size()
		newWidth := float32(currentSize.Width) * 1.2
		newHeight := float32(currentSize.Height) * 1.2
		mainWindow.Resize(fyne.NewSize(newWidth, newHeight))
	})
	minusButton := widget.NewButton("-", func() {
		currentSize := mainWindow.Canvas().Size()
		newWidth := float32(currentSize.Width) * 0.8
		newHeight := float32(currentSize.Height) * 0.8
		mainWindow.Resize(fyne.NewSize(newWidth, newHeight))
	})

	sizeMessage := widget.NewLabel("Change size")
	spacer := widget.NewLabel("")
	plusMinusContainer := container.NewHBox(spacer, minusButton, plusButton)
	sizePopUpContainer := container.NewVBox(sizeMessage, plusMinusContainer)

	sizePopUp := widget.NewPopUp(sizePopUpContainer, mainWindow.Canvas())

	sizeOption := fyne.NewMenuItem("Size", func() {
		sizePopUp.Show()
	})

	settingsMenu := fyne.NewMenu("Setting", themeOption, sizeOption)
	settingsMenuPopUp := widget.NewPopUpMenu(settingsMenu, mainWindow.Canvas())

	verticalToolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			searchMenuPopUp.Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			savePopUp.Resize(fyne.NewSize(200, 100))
			savePopUp.Show()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			settingsMenuPopUp.Show()
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
			return widget.NewLabel("Key options")
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

// EventListView creates and returns a Fyne List widget that displays the names of the given events.
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

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
