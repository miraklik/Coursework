package main

import (
	"bookstore/db"
	"bookstore/handlers"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	return db.ConnectDB()
}

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Bookstore")
	myWindow.Resize(fyne.NewSize(600, 500))

	tabs := container.NewAppTabs(
		container.NewTabItem("Управление книгами", createBooksTab(myWindow, db)),
		container.NewTabItem("Управление авторами", createAuthorsTab(myWindow, db)),
		container.NewTabItem("Поиск", createSearchTab(myWindow, db)),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func createBooksTab(window fyne.Window, db *gorm.DB) fyne.CanvasObject {
	heading := canvas.NewText("Управление книгами", theme.PrimaryColor())
	heading.TextSize = 24
	heading.Alignment = fyne.TextAlignCenter

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Введите название книги")

	authorEntry := widget.NewEntry()
	authorEntry.SetPlaceHolder("Введите автора книги")

	createCard := widget.NewCard("Добавить новую книгу", "", nil)
	createContent := container.NewVBox(
		widget.NewLabel("Название:"),
		titleEntry,
		widget.NewLabel("Автор:"),
		authorEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Создать", theme.ContentAddIcon(), func() {
				if titleEntry.Text == "" || authorEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("Необходимо заполнить все поля"), window)
					return
				}

				book, err := handlers.CreateBook(db, titleEntry.Text, authorEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Ошибка при создании книги: %v", err), window)
					return
				}

				dialog.ShowInformation("Успешно", fmt.Sprintf("Книга '%s' успешно добавлена", book.Title), window)
				titleEntry.SetText("")
				authorEntry.SetText("")
			}),
		),
	)
	createCard.SetContent(createContent)

	titleForDeleteBookEntry := widget.NewEntry()
	titleForDeleteBookEntry.SetPlaceHolder("Введите название книги для удаления")

	deleteCard := widget.NewCard("Удалить книгу", "", nil)
	deleteContent := container.NewVBox(
		widget.NewLabel("Название книги:"),
		titleForDeleteBookEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
				if titleForDeleteBookEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("Введите название книги"), window)
					return
				}

				confirmDialog := dialog.NewConfirm(
					"Подтверждение удаления",
					fmt.Sprintf("Вы уверены, что хотите удалить книгу '%s'?", titleForDeleteBookEntry.Text),
					func(ok bool) {
						if ok {
							err := handlers.DeleteBook(db, titleForDeleteBookEntry.Text)
							if err != nil {
								dialog.ShowError(fmt.Errorf("Ошибка при удалении книги: %v", err), window)
								return
							}

							dialog.ShowInformation("Успешно", "Книга успешно удалена", window)
							titleForDeleteBookEntry.SetText("")
						}
					},
					window,
				)
				confirmDialog.Show()
			}),
		),
	)
	deleteCard.SetContent(deleteContent)

	viewAllBooksCard := widget.NewCard("Просмотр всех книг", "", nil)
	booksList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Заголовок книги - Автор")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
		},
	)

	scrollContainer := container.NewVScroll(booksList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	viewAllBooksContent := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Загрузить все книги", theme.ViewRefreshIcon(), func() {
				books, err := handlers.GetAllBooks(db)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Ошибка при получении книг: %v", err), window)
					return
				}

				localBooks := *books

				booksList.Length = func() int { return len(*books) }
				booksList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
					label := obj.(*widget.Label)
					label.SetText(fmt.Sprintf("%s - %s", localBooks[id].Title, localBooks[id].Author))
				}
				booksList.Refresh()
			}),
			layout.NewSpacer(),
		),
		scrollContainer,
	)
	viewAllBooksCard.SetContent(viewAllBooksContent)

	// Компоновка вкладки
	return container.NewVBox(
		heading,
		container.NewPadded(
			container.NewVBox(
				createCard,
				widget.NewSeparator(),
				deleteCard,
				widget.NewSeparator(),
				viewAllBooksCard,
			),
		),
	)
}

func createAuthorsTab(window fyne.Window, db *gorm.DB) fyne.CanvasObject {
	heading := canvas.NewText("Управление авторами", theme.PrimaryColor())
	heading.TextSize = 24
	heading.Alignment = fyne.TextAlignCenter

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите имя автора")

	createAuthorCard := widget.NewCard("Добавить нового автора", "", nil)
	createAuthorContent := container.NewVBox(
		widget.NewLabel("Имя автора:"),
		nameEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Создать", theme.ContentAddIcon(), func() {
				if nameEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("Введите имя автора"), window)
					return
				}

				author, err := handlers.CreateAuthor(db, nameEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Ошибка при создании автора: %v", err), window)
					return
				}

				dialog.ShowInformation("Успешно", fmt.Sprintf("Автор '%s' успешно добавлен", author.Name), window)
				nameEntry.SetText("")
			}),
		),
	)
	createAuthorCard.SetContent(createAuthorContent)

	return container.NewVBox(
		heading,
		container.NewPadded(createAuthorCard),
	)
}

func createSearchTab(window fyne.Window, db *gorm.DB) fyne.CanvasObject {
	heading := canvas.NewText("Поиск книг", theme.PrimaryColor())
	heading.TextSize = 24
	heading.Alignment = fyne.TextAlignCenter

	titleSearchEntry := widget.NewEntry()
	titleSearchEntry.SetPlaceHolder("Введите название книги")

	authorSearchEntry := widget.NewEntry()
	authorSearchEntry.SetPlaceHolder("Введите имя автора")

	resultsCard := widget.NewCard("Результаты поиска", "", nil)

	resultsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Заголовок книги - Автор")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
		},
	)

	scrollContainer := container.NewVScroll(resultsList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	resultsContent := container.NewVBox(
		widget.NewLabel("Найденные книги:"),
		scrollContainer,
	)
	resultsCard.SetContent(resultsContent)

	updateSearchResults := func(books []interface{}) {
		resultsList.Length = func() int { return len(books) }
		resultsList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			book := books[id].(map[string]interface{})
			label.SetText(fmt.Sprintf("%v - %v", book["Title"], book["Author"]))
		}
		resultsList.Refresh()
	}

	searchByTitleCard := widget.NewCard("Поиск по названию", "", nil)
	searchByTitleContent := container.NewVBox(
		titleSearchEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Найти", theme.SearchIcon(), func() {
				if titleSearchEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("Введите название книги"), window)
					return
				}

				books, err := handlers.GetBookByAuthor(db, authorSearchEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Ошибка при поиске книг: %v", err), window)
					return
				}

				if len(*books) == 0 {
					dialog.ShowInformation("Результаты поиска", "Книги не найдены", window)
					return
				}

				var bookList []interface{}
				for _, book := range *books {
					bookMap := map[string]interface{}{
						"Title":  book.Title,
						"Author": book.Author,
					}
					bookList = append(bookList, bookMap)
				}

				updateSearchResults(bookList)
			}),
		),
	)
	searchByTitleCard.SetContent(searchByTitleContent)

	searchByAuthorCard := widget.NewCard("Поиск по автору", "", nil)
	searchByAuthorContent := container.NewVBox(
		authorSearchEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Найти", theme.SearchIcon(), func() {
				if authorSearchEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("Введите имя автора"), window)
					return
				}

				books, err := handlers.GetBookByAuthor(db, authorSearchEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Ошибка при поиске книг: %v", err), window)
					return
				}

				if len(*books) == 0 {
					dialog.ShowInformation("Результаты поиска", "Книги не найдены", window)
					return
				}

				var bookList []interface{}
				for _, book := range *books {
					bookMap := map[string]interface{}{
						"Title":  book.Title,
						"Author": book.Author,
					}
					bookList = append(bookList, bookMap)
				}

				updateSearchResults(bookList)
			}),
		),
	)
	searchByAuthorCard.SetContent(searchByAuthorContent)

	return container.NewVBox(
		heading,
		container.NewPadded(
			container.NewVBox(
				container.NewHSplit(
					searchByTitleCard,
					searchByAuthorCard,
				),
				resultsCard,
			),
		),
	)
}
