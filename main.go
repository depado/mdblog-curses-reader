package main

import (
	"fmt"

	gc "github.com/rthornton128/goncurses"
)

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	userArticles, err := fetchAllArticles()
	perror(err)

	// Initialize the standard screen
	stdscr, _ := gc.Init()
	defer gc.End()
	h, w := stdscr.MaxYX()

	// Initialize the library
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

	// Build the user menu items
	userItems := make([]*gc.MenuItem, len(userArticles))
	i := 0
	for val := range userArticles {
		userItems[i], _ = gc.NewItem(val, "")
		defer userItems[i].Free()
		i++
	}

	// Build the first user's article items
	articleItems := make([]*gc.MenuItem, len(userArticles[userItems[0].Name()]))
	i = 0
	for _, val := range userArticles[userItems[0].Name()] {
		articleItems[i], _ = gc.NewItem(val.title, "")
		defer articleItems[i].Free()
		i++
	}

	// Create the user menu
	userMenu, userMenuWin := newFeatureMenu(userItems, h, w/2, 0, 0)
	userMenu.Post()
	defer userMenu.UnPost()
	userMenuWin.Refresh()

	// Create the article menu
	articleMenu, articleMenuWin := newFeatureMenu(articleItems, h/2, w/2, 0, w/2)
	articleMenu.Post()
	defer articleMenu.UnPost()
	articleMenuWin.Refresh()

	aInfoWin, _ := gc.NewWindow(h/2+1, w/2, h/2, w/2+1)
	displayInfo(aInfoWin, 0, 0, userItems, userArticles)

	active := 0
	activeA := 0
	currentMenuWin := userMenuWin
	currentMenu := userMenu
	for {
		gc.Update()
		ch := currentMenuWin.GetChar()
		switch ch {
		case 'q':
			return
		case gc.KEY_UP:
			if currentMenu == userMenu {
				if active != 0 {
					active--
					articleItems = make([]*gc.MenuItem, len(userArticles[userItems[active].Name()]))
					i = 0
					for _, val := range userArticles[userItems[active].Name()] {
						articleItems[i], _ = gc.NewItem(val.title, "")
						defer articleItems[i].Free()
						i++
					}
					articleMenu.UnPost()
					articleMenu.SetItems(articleItems)
					articleMenu.Post()
					articleMenuWin.Refresh()
					displayInfo(aInfoWin, active, activeA, userItems, userArticles)
				}
			} else {
				if activeA != 0 {
					activeA--
					displayInfo(aInfoWin, active, activeA, userItems, userArticles)
				}
			}
			currentMenu.Driver(gc.DriverActions[ch])
		case gc.KEY_DOWN:
			if currentMenu == userMenu {
				if active != len(userItems)-1 {
					active++
					articleItems = make([]*gc.MenuItem, len(userArticles[userItems[active].Name()]))
					i = 0
					for _, val := range userArticles[userItems[active].Name()] {
						articleItems[i], _ = gc.NewItem(val.title, "")
						defer articleItems[i].Free()
						i++
					}
					articleMenu.UnPost()
					articleMenu.SetItems(articleItems)
					articleMenu.Post()
					articleMenuWin.Refresh()
					displayInfo(aInfoWin, active, activeA, userItems, userArticles)
				}
			} else {
				if activeA != len(articleItems)-1 {
					activeA++
					aInfoWin.MovePrint(1, 1, "Auteur : ", userItems[active].Name())
					aInfoWin.Refresh()
					displayInfo(aInfoWin, active, activeA, userItems, userArticles)
				}
			}
			currentMenu.Driver(gc.DriverActions[ch])
		case gc.KEY_LEFT:
			currentMenu = userMenu
			currentMenuWin = userMenuWin
			activeA = 0
		case gc.KEY_RIGHT:
			currentMenu = articleMenu
			currentMenuWin = articleMenuWin
		case gc.KEY_ENTER, gc.KEY_RETURN, gc.Key('\r'):
			if currentMenu == userMenu {
				currentMenu = articleMenu
				currentMenuWin = articleMenuWin
				activeA = 0
			} else {
				gc.End()
				content := dlContent(userArticles[userItems[active].Name()][activeA].url + "/ansi")
				fmt.Println(content)
				return
			}
		default:
			currentMenu.Driver(gc.DriverActions[ch])
		}
	}
}
