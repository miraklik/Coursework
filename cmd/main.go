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
		container.NewTabItem("Управление студентами", createStudentsTab(myWindow, db)),
		container.NewTabItem("Управление организациями", createOrganizationTab(myWindow, db)),
		container.NewTabItem("Поиск", createSearchTab(myWindow, db)),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func createStudentsTab(window fyne.Window, db *gorm.DB) fyne.CanvasObject {
	heading := canvas.NewText("Управление Студентами", theme.PrimaryColor())
	heading.TextSize = 24
	heading.Alignment = fyne.TextAlignCenter

	NameEntry := widget.NewEntry()
	NameEntry.SetPlaceHolder("Введите имя студента")

	GroupEntry := widget.NewEntry()
	GroupEntry.SetPlaceHolder("Введите название группы")

	createCard := widget.NewCard("Добавить нового студента", "", nil)
	createContent := container.NewVBox(
		widget.NewLabel("Имя:"),
		NameEntry,
		widget.NewLabel("Группа:"),
		GroupEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Создать", theme.ContentAddIcon(), func() {
				if NameEntry.Text == "" || GroupEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("необходимо заполнить все поля"), window)
					return
				}

				student, err := handlers.CreateStudents(db, NameEntry.Text, GroupEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка при добавлении студента: %v", err), window)
					return
				}

				dialog.ShowInformation("Успешно", fmt.Sprintf("Студент '%s' успешно добавлен", student.Full_name), window)
				NameEntry.SetText("")
				GroupEntry.SetText("")
			}),
		),
	)
	createCard.SetContent(createContent)

	NameForDeleteBookEntry := widget.NewEntry()
	NameForDeleteBookEntry.SetPlaceHolder("Введите имя студента для удаления")

	deleteCard := widget.NewCard("Удалить Студента", "", nil)
	deleteContent := container.NewVBox(
		widget.NewLabel("Имя Студента:"),
		NameForDeleteBookEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
				if NameForDeleteBookEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("введите имя студента"), window)
					return
				}

				confirmDialog := dialog.NewConfirm(
					"Подтверждение удаления",
					fmt.Sprintf("Вы уверены, что хотите удалить студента '%s'?", NameForDeleteBookEntry.Text),
					func(ok bool) {
						if ok {
							err := handlers.DeleteStudens(db, NameForDeleteBookEntry.Text)
							if err != nil {
								dialog.ShowError(fmt.Errorf("ошибка при удалении студента: %v", err), window)
								return
							}

							dialog.ShowInformation("Успешно", "Студент успешно удален", window)
							NameForDeleteBookEntry.SetText("")
						}
					},
					window,
				)
				confirmDialog.Show()
			}),
		),
	)
	deleteCard.SetContent(deleteContent)

	viewAllBooksCard := widget.NewCard("Просмотр всех студентов", "", nil)
	booksList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Имя студента - Группа")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
		},
	)

	scrollContainer := container.NewVScroll(booksList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	viewAllBooksContent := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Загрузить всех студентов", theme.ViewRefreshIcon(), func() {
				students, err := handlers.GetAllStudents(db)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка при получении студентов: %v", err), window)
					return
				}

				localStudents := *students

				booksList.Length = func() int { return len(*students) }
				booksList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
					label := obj.(*widget.Label)
					label.SetText(fmt.Sprintf("%s - %s", localStudents[id].Full_name, localStudents[id].Group_name))
				}
				booksList.Refresh()
			}),
			layout.NewSpacer(),
		),
		scrollContainer,
	)
	viewAllBooksCard.SetContent(viewAllBooksContent)

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

func createOrganizationTab(window fyne.Window, db *gorm.DB) fyne.CanvasObject {
	heading := canvas.NewText("Управление организациями", theme.PrimaryColor())
	heading.TextSize = 24
	heading.Alignment = fyne.TextAlignCenter

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите название организации")

	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Введите адрес организации")

	createCard := widget.NewCard("Добавить новую организацию", "", nil)
	createContent := container.NewVBox(
		widget.NewLabel("Название организации:"),
		nameEntry,
		widget.NewLabel("Адрес организации:"),
		addressEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Создать", theme.ContentAddIcon(), func() {
				if nameEntry.Text == "" || addressEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("необходимо заполнить все поля"), window)
					return
				}

				organization, err := handlers.CreateOrganization(db, nameEntry.Text, addressEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка при создании организации: %v", err), window)
					return
				}

				dialog.ShowInformation("Успешно", fmt.Sprintf("Организация '%s' успешно добавлена", organization.Name), window)
				nameEntry.SetText("")
				addressEntry.SetText("")
			}),
		),
	)
	createCard.SetContent(createContent)

	nameForDeleteEntry := widget.NewEntry()
	nameForDeleteEntry.SetPlaceHolder("Введите название организации для удаления")

	deleteCard := widget.NewCard("Удалить организацию", "", nil)
	deleteContent := container.NewVBox(
		widget.NewLabel("Название организации:"),
		nameForDeleteEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
				if nameForDeleteEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("введите название организации"), window)
					return
				}

				confirmDialog := dialog.NewConfirm(
					"Подтверждение удаления",
					fmt.Sprintf("Вы уверены, что хотите удалить организацию '%s'?", nameForDeleteEntry.Text),
					func(ok bool) {
						if ok {
							err := handlers.DeleteOrganization(db, nameForDeleteEntry.Text)
							if err != nil {
								dialog.ShowError(fmt.Errorf("ошибка при удалении организации: %v", err), window)
								return
							}

							dialog.ShowInformation("Успешно", "Организация успешно удалена", window)
							nameForDeleteEntry.SetText("")
						}
					},
					window,
				)
				confirmDialog.Show()
			}),
		),
	)
	deleteCard.SetContent(deleteContent)

	viewAllOrganizationsCard := widget.NewCard("Просмотр всех организаций", "", nil)
	organizationsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Название организации - Адрес")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
		},
	)

	scrollContainer := container.NewVScroll(organizationsList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	viewAllOrganizationsContent := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Загрузить все организации", theme.ViewRefreshIcon(), func() {
				organizations, err := handlers.GetAllOrganizations(db)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка при получении организаций: %v", err), window)
					return
				}

				localOrganizations := *organizations

				organizationsList.Length = func() int { return len(*organizations) }
				organizationsList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
					label := obj.(*widget.Label)
					label.SetText(fmt.Sprintf("%s - %s", localOrganizations[id].Name, localOrganizations[id].Address))
				}
				organizationsList.Refresh()
			}),
			layout.NewSpacer(),
		),
		scrollContainer,
	)
	viewAllOrganizationsCard.SetContent(viewAllOrganizationsContent)

	return container.NewVBox(
		heading,
		container.NewPadded(
			container.NewVBox(
				createCard,
				widget.NewSeparator(),
				deleteCard,
				widget.NewSeparator(),
				viewAllOrganizationsCard,
			),
		),
	)
}

func createSearchTab(window fyne.Window, db *gorm.DB) fyne.CanvasObject {
	heading := canvas.NewText("Поиск", theme.PrimaryColor())
	heading.TextSize = 24
	heading.Alignment = fyne.TextAlignCenter

	NameSearchEntry := widget.NewEntry()
	NameSearchEntry.SetPlaceHolder("Введите имя студента")

	OrganizationSearchEntry := widget.NewEntry()
	OrganizationSearchEntry.SetPlaceHolder("Введите название организации")

	resultsCard := widget.NewCard("Результаты поиска", "", nil)

	resultsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Результат поиска")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
		},
	)

	scrollContainer := container.NewVScroll(resultsList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	resultsContent := container.NewVBox(
		widget.NewLabel("Найденные результаты:"),
		scrollContainer,
	)
	resultsCard.SetContent(resultsContent)

	updateSearchResults := func(items []string) {
		resultsList.Length = func() int { return len(items) }
		resultsList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText(items[id])
		}
		resultsList.Refresh()
	}

	searchByTitleCard := widget.NewCard("Поиск студента по имени", "", nil)
	searchByTitleContent := container.NewVBox(
		NameSearchEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Найти", theme.SearchIcon(), func() {
				if NameSearchEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("введите имя студента"), window)
					return
				}

				student, err := handlers.GetStudentByName(db, NameSearchEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка при поиске студента: %v", err), window)
					return
				}

				if student == nil {
					dialog.ShowInformation("Результаты поиска", "Студент не найден", window)
					return
				}

				var studentResults []string
				studentResults = append(studentResults, fmt.Sprintf("Студент: %s - Группа: %s", (*student)[0].Full_name, (*student)[0].Group_name))

				updateSearchResults(studentResults)
			}),
		),
	)
	searchByTitleCard.SetContent(searchByTitleContent)

	searchByAuthorCard := widget.NewCard("Поиск по названию организации", "", nil)
	searchByAuthorContent := container.NewVBox(
		OrganizationSearchEntry,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Найти", theme.SearchIcon(), func() {
				if OrganizationSearchEntry.Text == "" {
					dialog.ShowError(fmt.Errorf("введите название организации"), window)
					return
				}

				organization, err := handlers.GetOrganizationsByName(db, OrganizationSearchEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка при поиске организации: %v", err), window)
					return
				}

				if organization == nil {
					dialog.ShowInformation("Результаты поиска", "Организация не найдена", window)
					return
				}

				var organizationResults []string
				organizationResults = append(organizationResults, fmt.Sprintf("Организация: %s - Адрес: %s", (*organization)[0].Name, (*organization)[0].Address))

				updateSearchResults(organizationResults)
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
