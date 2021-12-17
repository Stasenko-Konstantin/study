package src

import (
	"fmt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"net"
	"os"
	"strings"
	_ "strings"
)

var (
	mylog *myLogger
	w     fyne.Window
)

var (
	clients   = table{row{"ID", "ФИО", "Место жительства", ""}}
	cassettes = table{row{"ID", "Стоимость", "Фильм", ""}}
	films     = table{row{"Название", "Год выпуска", "Режиссер", "Жанр",
		"Хронометраж", "Производитель", ""}}
	librarians = table{row{"ID", "ФИО", ""}}
	givings    = table{row{"ID", "Клиент", "Кассета", "Выдана"}}
)

func Start() {
	mylog = newLogger()

	defer func() {
		if r := recover(); r != nil {
			s := r.(error)
			mylog.Write([]byte(s.Error() + " aga\n"))
		}
	}()

	a := app.New()
	w = a.NewWindow("практика")
	w.Resize(fyne.NewSize(1800, 900))

	connect()
	mylog.Write([]byte("DB connection complete\n"))

	clientsT := widget.NewTable(
		func() (int, int) {
			return len(clients), len(clients[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0] = widget.NewLabel(clients[i.Row][i.Col])
		})
	clientsT.SetColumnWidth(0, 50.0)
	clientsT.SetColumnWidth(1, 300.0)
	clientsT.SetColumnWidth(2, 650.0)

	cassettesT := widget.NewTable(
		func() (int, int) {
			return len(cassettes), len(cassettes[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0] = widget.NewLabel(cassettes[i.Row][i.Col])
		})
	cassettesT.SetColumnWidth(0, 50.0)
	cassettesT.SetColumnWidth(1, 100.0)
	cassettesT.SetColumnWidth(2, 230.0)

	filmsT := widget.NewTable(
		func() (int, int) {
			return len(films), len(films[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0] = widget.NewLabel(films[i.Row][i.Col])
		})
	filmsT.SetColumnWidth(0, 200.0)
	filmsT.SetColumnWidth(1, 130.0)
	filmsT.SetColumnWidth(2, 200.0)
	filmsT.SetColumnWidth(3, 150.0)
	filmsT.SetColumnWidth(4, 150.0)
	filmsT.SetColumnWidth(5, 200.0)

	librariansT := widget.NewTable(
		func() (int, int) {
			return len(librarians), len(librarians[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0] = widget.NewLabel(librarians[i.Row][i.Col])
		})
	librariansT.SetColumnWidth(0, 50.0)
	librariansT.SetColumnWidth(1, 300.0)

	givingsT := widget.NewTable(
		func() (int, int) {
			return len(givings), len(givings[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0] = widget.NewLabel(givings[i.Row][i.Col])
		})
	givingsT.SetColumnWidth(0, 50.0)
	givingsT.SetColumnWidth(1, 100.0)
	givingsT.SetColumnWidth(2, 100.0)
	givingsT.SetColumnWidth(3, 100.0)

	mylog.Write([]byte("Table template ready\n"))

	checkErr := func(err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	result := widget.NewMultiLineEntry()
	result.Disable()
	result.Move(fyne.NewPos(0.0, 40.0))
	result.Resize(fyne.NewSize(1800, 250))
	requestE := widget.NewEntry()

	//1 выбрать фильмы по цене
	var priceS *widget.Select
	ops := []string{">", ">=", "<", "<=", "=="}
	op := ops[0]
	opS := widget.NewSelect(ops, func(opfS string) {
		op = opfS
		priceS.SetSelectedIndex(priceS.SelectedIndex())
	})
	prices := makePrices()
	priceS = widget.NewSelect(prices, func(price string) {
		requestE.SetText("select cs.id, f.name, f.year, f.director, cs.price from cassettes as cs join films as f on cs.film = f.name and cs.year = f.year where cs.price " + op + " " +
			price)
	})
	r1 := container.NewVBox(container.NewHBox(opS), priceS)
	priceS.Move(fyne.NewPos(0, 10))
	priceS.Refresh()
	opS.SetSelectedIndex(0)

	//2 выбрать фильмы по цене
	var priceS *widget.Select
	ops := []string{">", ">=", "<", "<=", "=="}
	op := ops[0]
	opS := widget.NewSelect(ops, func(opfS string) {
		op = opfS
		priceS.SetSelectedIndex(priceS.SelectedIndex())
	})
	prices := makePrices()
	priceS = widget.NewSelect(prices, func(price string) {
		requestE.SetText("select cs.id, f.name, f.year, f.director, cs.price from cassettes as cs join films as f on cs.film = f.name and cs.year = f.year where cs.price " + op + " " +
			price)
	})
	r1 := container.NewVBox(container.NewHBox(opS), priceS)
	priceS.Move(fyne.NewPos(0, 10))
	priceS.Refresh()
	opS.SetSelectedIndex(0)

	priceS.SetSelectedIndex(0) // из первого запроса

	hideAll := func() {
		r1.Hide()
		r2.Hide()
		r3.Hide()
		r4.Hide()
		r5.Hide()
		r6.Hide()
		r7.Hide()
		r8.Hide()
	}
	hideAll()
	r1.Show()

	rs := []string{
		"Фильмы по цене",
		"Фильмы по хронометражу",
		"Фильмы по жанрам",
		"Фильмы по годам",
		"Фильмы по режиссерам",
		"(Не-)выданные кассеты",
		"Кассеты по библиотекарям",
	}
	requests := widget.NewSelect(rs, func(r string) {
		hideAll()
		switch r {
		case rs[0]:
			r1.Show()
			priceS.SetSelectedIndex(0)
		case rs[1]:
			r2.Show()
			timeS.SetSelectedIndex(0)
		case rs[2]:
			r3.Show()
			genreS.SetSelectedIndex(0)
		case rs[3]:
			r4.Show()
			yearS.SetSelectedIndex(0)
		case rs[4]:
			r5.Show()
			dirS.SetSelectedIndex(0)
		case rs[5]:
			r6.Show()
			casS.SetSelectedIndex(0)
		case rs[6]:
			r7.Show()
			libS.SetSelectedIndex(0)
		}
	})
	requests.SetSelectedIndex(0)

	requestsC := container.NewVBox(
		requests,
		r1,
		r2,
		r3,
		r4,
		r5,
		r6,
		r7,
		requestE,
	)
	requestsC.Move(fyne.NewPos(0.0, 300.0))
	requestsC.Resize(fyne.NewSize(1600, 200))
	requestsC.Refresh()

	mylog.Write([]byte("Request fields ready\n"))

	tabs := container.NewAppTabs(
		container.NewTabItem("Клиенты", clientsT),
		container.NewTabItem("Кассеты", cassettesT),
		container.NewTabItem("Фильмы", filmsT),
		container.NewTabItem("Библиотекари", librariansT),
		container.NewTabItem("Выдачи", givingsT),
		container.NewTabItem("Запросы", container.NewWithoutLayout(
			container.NewHBox(
				widget.NewButton("Запрос", func() {
					go func() {
						pc, err := net.ListenPacket("udp", myIP.String()+":12345")
						if err != nil {
							mylog.Write([]byte(err.Error()))
						}
						send(allIP(myIP.String()), encode(requestE.Text+";"), mylog)
						buf := make([]byte, 10000)
						_, _, err = pc.ReadFrom(buf)
						if err != nil {
							mylog.Write([]byte("AGAAA" + err.Error()))
							return
						}
						msg := decode(strings.Split(string(buf), "$|")[1])
						result.SetText(read(msg))
					}()
				}),
				widget.NewButton("Печать", func() {
					dialog.ShowFileSave(func(closer fyne.URIWriteCloser, err error) {
						if closer == nil {
							return
						}
						if err != nil {
							panic(err)
						}
						path := closer.URI().Path()
						closer.Close()
						text := []byte(result.Text)
						err = os.WriteFile(path, text, 0644)
						checkErr(err)
					}, w)
				}),
			),
			result,
			requestsC,
		)))

	prog := "Программа по учебной практике №10 \"Разработка программы для работы с удаленной базой данных\""
	author := "Курсант 432 гр. ТАТК ГА Стасенко Константин"

	mainMenu := fyne.NewMainMenu(fyne.NewMenu("Меню",
		fyne.NewMenuItem("О программе", func() { dialog.ShowInformation("О программе", prog, w) }),
		fyne.NewMenuItem("Об авторе", func() { dialog.ShowInformation("Об авторе", author, w) }),
	))

	w.SetMainMenu(mainMenu)

	mylog.Write([]byte("Start render\n"))

	w.SetContent(tabs)
	w.ShowAndRun()
}
