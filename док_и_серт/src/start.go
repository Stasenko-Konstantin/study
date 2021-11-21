package src

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
	"time"
)

func sum(smth string) (int, string) {
	var r int
	for _, e := range smth {
		e, err := strconv.Atoi(string(e))
		if err != nil {
			return 0, err.Error()
		}
		r += e
	}
	return r, "nil"
}

func Start() {
	a := app.New()
	w := a.NewWindow("практика")
	w.Resize(fyne.NewSize(400, 200))
	w.CenterOnScreen()

	dayE := widget.NewEntry()
	monthE := widget.NewEntry()
	happyYearL := widget.NewLabel("")

	w.SetContent(container.NewVBox(
		widget.NewLabel("Введите день и месяц своего рождения\nПосле нажатия кнопки \"Ok\" выведется ваш\nближайший счастливый год"),
		container.NewHBox(
			widget.NewLabel("День:"),
			dayE,
			widget.NewLabel("Месяц:"),
			monthE,
			widget.NewButton("Ok", func() {
				day, err := strconv.Atoi(dayE.Text)
				if err != nil {
					dialog.ShowError(errors.New("День должен быть числом!"), w)
					return
				}
				month, err := strconv.Atoi(monthE.Text)
				if err != nil {
					dialog.ShowError(errors.New("Месяц должен быть числом!"), w)
					return
				}
				if month < 1 || month > 12 {
					dialog.ShowError(errors.New("Месяц должен быть от 1 до 12!"), w)
					return
				}
				limit := 31
				if (month < 8 && month % 2 == 0) || (month > 8 && month % 2 == 1) {
					limit = 30
				}
				if month == 2 {
					limit = 28
				}
				if day < 1 || day > limit {
					dialog.ShowError(errors.New("День должен быть от 1 до " + strconv.Itoa(limit)), w)
					return
				}
				sumD, myerr := sum(strconv.Itoa(day))
				if myerr != "nil" {
					dialog.ShowError(errors.New(myerr), w)
					return
				}
				sumM, myerr := sum(strconv.Itoa(month))
				if myerr != "nil" {
					dialog.ShowError(errors.New(myerr), w)
					return
				}
				currYear := strings.Split(time.Now().String(), "-")[0]
				for {
					fmt.Println(currYear)
					sumY, myerr := sum(currYear)
					if myerr != "nil" {
						dialog.ShowError(errors.New(myerr), w)
						return
					}
					if sumY % 10 == sumD % 10 || sumY % 10 == sumM % 10 {
						happyYearL.SetText(currYear)
						return
					}
					prevYear, err := strconv.Atoi(currYear)
					if err != nil {
						dialog.ShowError(err, w)
						return
					}
					currYear = strconv.Itoa(prevYear + 1)
				}
			}),
		),
		container.NewHBox(
			widget.NewLabel("Ближайший счастливый год: "),
			happyYearL,
		),
	))

	prog := "Программа по \"Документации и сертификации\":\n\"Разработка программы для вычисления\nсвоего ближайшего счастливого года\""
	author := "Курсант 432 гр. ТАТК ГА Стасенко Константин"
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("Меню",
		fyne.NewMenuItem("О программе", func() { dialog.ShowInformation("О программе", prog, w) }),
		fyne.NewMenuItem("Об авторе", func() { dialog.ShowInformation("Об авторе", author, w) }),
	))
	w.SetMainMenu(mainMenu)

	w.ShowAndRun()
}
