package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
	"strings"
	"time"
)

type book []string

var books = [][]string{
	book{"АВТОР", "КНИГА", "ИЗДАТЕЛЬСТВО", "ГОД ИЗДАНИЯ"},
}

func checkErr(err error) {
	fmt.Println(err.Error())
}

func checkDate(olddate string, w *fyne.Window) {
	date, err := strconv.Atoi(olddate)
	if err != nil {
		dialog.ShowError(errors.New("Неверная дата - "+olddate+", вводите только цифры!"), *w)
	}
	curr := time.Now().Year()
	if date < 1700 || date > curr+1 {
		dialog.ShowError(errors.New("Неверная дата - "+olddate+", вводите дату в диапазоне от 1700 до "+strconv.Itoa(curr+1)), *w)
	}
}

func getPeriod(start, end string) string {
	var res string
	for _, b := range books {
		year, _ := strconv.Atoi(b[3])
		start, _ := strconv.Atoi(start)
		end, _ := strconv.Atoi(end)
		if year >= start && year <= end {
			res += stringConcat(b, ", ") + "\n"
		}
	}
	return res
}

func stringConcat(vals []string, del string) string {
	var r string
	for _, e := range vals {
		r += e + " " + del
	}
	return r
}

func Start() {

	defer func() {
		if r := recover(); r != nil {
			s := r.(error)
			fmt.Println(s.Error() + " aga")
		}
	}()

	text, err := os.ReadFile("test.txt")
	if err != nil {
		checkErr(err)
	}
	texts := strings.Split(string(text), "\n")
	for _, t := range texts {
		var r []string
		elems := strings.Split(t, ", ")
		if len(elems) != 4 {
			continue
		}
		for _, e := range elems {
			r = append(r, e)
		}
		books = append(books, r)
	}

	a := app.New()
	w := a.NewWindow("модуль")
	w.Resize(fyne.NewSize(1800, 900))

	booksT := widget.NewTable(
		func() (int, int) {
			return len(books), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			book := books[i.Row][i.Col]
			o.(*widget.Label).SetText(book)
		})
	booksT.SetColumnWidth(0, 170.0)
	booksT.SetColumnWidth(1, 300.0)
	booksT.SetColumnWidth(2, 170.0)
	booksT.SetColumnWidth(3, 100.0)

	result := widget.NewMultiLineEntry()
	result.Disable()
	result.Move(fyne.NewPos(0.0, 0.0))
	result.Resize(fyne.NewSize(1800, 250))
	start := widget.NewEntry()
	end := widget.NewEntry()

	requestC := container.NewVBox(
		container.NewVBox(
			widget.NewLabel("Дата начала:"),
			start,
			widget.NewLabel("Дата окончания:"),
			end,
		),
		container.NewHBox(
			widget.NewButton("Запрос", func() {
				go func() {
					start := start.Text
					end := end.Text
					checkDate(start, &w)
					checkDate(end, &w)
					result.SetText(getPeriod(start, end))
				}()
			}),
			widget.NewButton("Печать", func() {
				dialog.ShowFileSave(func(closer fyne.URIWriteCloser, err error) {
					defer func() {
						if r := recover(); r != nil {
							s := r.(error)
							fmt.Println(s.Error() + " aga")
						}
					}()
					if closer == nil {
						return
					}
					if err != nil {
						panic(err)
					}
					path := closer.URI().Path()
					closer.Close()
					text := []byte(result.Text)
					os.WriteFile(path, text, 0644)
				}, w)
			}),
		),
	)
	requestC.Move(fyne.NewPos(0.0, 255.0))
	requestC.Resize(fyne.NewSize(200, 200))
	requestC.Refresh()

	tabs := container.NewAppTabs(
		container.NewTabItem("Книги", booksT),
		container.NewTabItem("Запросы", container.NewWithoutLayout(
			result,
			requestC,
		)))

	prog := "Система обслуживания читателя в библиотеке"
	author := "Курсант 432 гр. ТАТК ГА Стасенко Константин"

	mainMenu := fyne.NewMainMenu(fyne.NewMenu("Меню",
		fyne.NewMenuItem("О программе", func() { dialog.ShowInformation("О программе", prog, w) }),
		fyne.NewMenuItem("Об авторе", func() { dialog.ShowInformation("Об авторе", author, w) }),
	))

	w.SetMainMenu(mainMenu)

	w.SetContent(tabs)
	w.ShowAndRun()
}
