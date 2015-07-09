package main

import (
	gc "github.com/rthornton128/goncurses"
)

// Returns a new Menu containing the items and a new Window containing the menu.
func newFeatureMenu(items []*gc.MenuItem, h, w, y, x int) (menu *gc.Menu, menuWin *gc.Window, err error) {
	if menu, err = gc.NewMenu(items); err != nil {
		return
	}
	if menuWin, err = gc.NewWindow(h, w, y, x); err != nil {
		return
	}
	menuWin.Keypad(true)
	menu.SetWindow(menuWin)
	menu.SubWindow(menuWin.Derived(h-2, w-2, 1, 1))
	menu.Format(h-2, 1)
	menu.Mark("")
	menuWin.Box(0, 0)
	return
}

// Function to call to update the information in the information Window.
func displayInfo(win *gc.Window, active, activeA int, userItems []*gc.MenuItem, userArticles map[string][]articleType) {
	win.Clear()
	win.MovePrint(1, 0, "Auteur : ", userItems[active].Name())
	win.MovePrint(2, 0, "Date : ", userArticles[userItems[active].Name()][activeA].date)
	win.MovePrint(4, 0, userArticles[userItems[active].Name()][activeA].url)
	win.Refresh()
}
