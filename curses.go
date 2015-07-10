package main

import gc "github.com/rthornton128/goncurses"

// Returns an articleType that can be downloaded and displayed.
func ncurses(articles map[string][]articleType) (article articleType, err error) {
	// Initialize the standard screen
	stdscr, err := gc.Init()
	if err != nil {
		return
	}
	defer gc.End()
	h, w := stdscr.MaxYX()

	// Initialize the library
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

	// Build the user menu items
	userItems := make([]*gc.MenuItem, len(articles))
	i := 0
	for val := range articles {
		userItems[i], err = gc.NewItem(val, "")
		if err != nil {
			return
		}
		defer userItems[i].Free()
		i++
	}

	// Build the first user's article items
	articleItems := make([]*gc.MenuItem, len(articles[userItems[0].Name()]))
	i = 0
	for _, val := range articles[userItems[0].Name()] {
		articleItems[i], err = gc.NewItem(val.title, "")
		if err != nil {
			return
		}
		defer articleItems[i].Free()
		i++
	}

	// Create the user menu
	userMenu, userMenuWin, err := newFeatureMenu(userItems, h, w/2, 0, 0)
	if err != nil {
		return
	}
	// Post the user menu and refresh the user window
	userMenu.Post()
	defer userMenu.UnPost()
	userMenuWin.Refresh()

	// Create the article menu
	articleMenu, articleMenuWin, err := newFeatureMenu(articleItems, h/2, w/2, 0, w/2)
	if err != nil {
		return
	}
	// Post the article menu and refresh the article window
	articleMenu.Post()
	defer articleMenu.UnPost()
	articleMenuWin.Refresh()

	// Create the information window used to display article information
	aInfoWin, err := gc.NewWindow(h/2+1, w/2, h/2, w/2+1)
	if err != nil {
		return
	}
	displayInfo(aInfoWin, 0, 0, userItems, articles)

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
					articleItems = make([]*gc.MenuItem, len(articles[userItems[active].Name()]))
					i = 0
					for _, val := range articles[userItems[active].Name()] {
						articleItems[i], _ = gc.NewItem(val.title, "")
						defer articleItems[i].Free()
						i++
					}
					articleMenu.UnPost()
					articleMenu.SetItems(articleItems)
					articleMenu.Post()
					articleMenuWin.Refresh()
					displayInfo(aInfoWin, active, activeA, userItems, articles)
				}
			} else {
				if activeA != 0 {
					activeA--
					displayInfo(aInfoWin, active, activeA, userItems, articles)
				}
			}
			currentMenu.Driver(gc.DriverActions[ch])
		case gc.KEY_DOWN:
			if currentMenu == userMenu {
				if active != len(userItems)-1 {
					active++
					articleItems = make([]*gc.MenuItem, len(articles[userItems[active].Name()]))
					i = 0
					for _, val := range articles[userItems[active].Name()] {
						articleItems[i], _ = gc.NewItem(val.title, "")
						defer articleItems[i].Free()
						i++
					}
					articleMenu.UnPost()
					articleMenu.SetItems(articleItems)
					articleMenu.Post()
					articleMenuWin.Refresh()
					displayInfo(aInfoWin, active, activeA, userItems, articles)
				}
			} else {
				if activeA != len(articleItems)-1 {
					activeA++
					aInfoWin.MovePrint(1, 1, "Auteur : ", userItems[active].Name())
					aInfoWin.Refresh()
					displayInfo(aInfoWin, active, activeA, userItems, articles)
				}
			}
			currentMenu.Driver(gc.DriverActions[ch])
		case gc.KEY_LEFT:
			// If the LEFT key is pressed, then set the active menu and active
			// menu window to the user ones. Also resets the current article counter
			if currentMenu == articleMenu {
				currentMenu = userMenu
				currentMenuWin = userMenuWin
				activeA = 0
			}
		case gc.KEY_RIGHT:
			// If the RIGHT key is pressed, then set the active menu and active
			// menu window to the article one.
			if currentMenu == userMenu {
				currentMenu = articleMenu
				currentMenuWin = articleMenuWin
			}
		case gc.KEY_ENTER, gc.KEY_RETURN, gc.Key('\r'):
			// If one of the Enter/Return key is pressed, depending on the current
			// menu, switch to the article menu or start downloading the selected
			// article.
			if currentMenu == userMenu {
				currentMenu = articleMenu
				currentMenuWin = articleMenuWin
				activeA = 0
			} else {
				article = articles[userItems[active].Name()][activeA]
				return
			}
		default:
			currentMenu.Driver(gc.DriverActions[ch])
		}
	}
}
