package src

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"os"
	_ "strings"
)

var (
	mylog *myLogger
	w     fyne.Window
)

var (
	patients = [][]string{[]string{"ID", "Страховая компания", "ФИО", "Место жительства",
		"Дата рождения", "Пол", "Участок", ""}}
	doctors = [][]string{[]string{"ФИО", "Отделение", "Специализация", "Участок", ""}}
	talons  = [][]string{[]string{"ID", "Дата и время приема", "Доктор", "Пациент", ""}}
)

func Start() {
	mylog = newLogger()
	var ch chan string
	go listen(ch)

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

	//for _, e := range dbs {
	//	switch e.(type) {
	//	case []patient:
	//		for _, p := range e.([]patient) {
	//			var sex string
	//			if p.sex {
	//				sex = "Ж."
	//			} else {
	//				sex = "М."
	//			}
	//			birth := strings.Split(p.birth, "T")[0]
	//			patients = append(patients, []string{strconv.Itoa(p.id), p.insuranceCompany, p.sfm, p.residence,
	//				birth, sex, strconv.Itoa(p.district), ""})
	//		}
	//	case []doctor:
	//		for _, d := range e.([]doctor) {
	//			doctors = append(doctors, []string{d.sfm, d.department, d.specialization, strconv.Itoa(int(d.district.Int32)), ""})
	//		}
	//	case []talon:
	//		for _, t := range e.([]talon) {
	//			r := strings.Split(t.reception, "T")
	//			r[1] = r[1][0:5]
	//			reception := stringConcat(r, " ")
	//			talons = append(talons, []string{strconv.Itoa(t.id), reception, t.doctor, t.patient, ""})
	//		}
	//	}
	//}

	mylog.Write([]byte("Data is full\n"))

	patientsT := widget.NewTable(
		func() (int, int) {
			return len(patients), len(patients[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			//if i.Row == 0 {
			//	o.(*fyne.Container).Objects[0] = widget.NewLabel(patients[i.Row][i.Col])
			//} else if i.Col == 7 {
			//	o.(*fyne.Container).Objects[0] = widget.NewButton("", func() {
			//		fmt.Println(i.Col, i.Row)
			//	})
			//} else {
			//	o.(*fyne.Container).Objects[0].(*widget.Entry).SetText(patients[i.Row][i.Col])
			//}
		})
	patientsT.SetColumnWidth(0, 50.0)

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
	var filmS *widget.Select
	ops := []string{">", ">=", "<", "<=", "=="}
	op := ops[0]
	opS := widget.NewSelect(ops, func(opfS string) {
		op = opfS
		filmS.SetSelectedIndex(filmS.SelectedIndex())
	})
	films := makePers()
	filmS = widget.NewSelect(films, func(film string) {
		requestE.SetText("select cs.id, f.name, f.year, f.director, cs.price from cassettes as cs join films as f on cs.film = f.name and cs.year = f.year where cs.price " + op + " " +
			film)
	})
	r1 := container.NewVBox(container.NewHBox(opS), filmS)
	filmS.Move(fyne.NewPos(0, 10))
	filmS.Refresh()
	opS.SetSelectedIndex(0)

	filmS.SetSelectedIndex(0) // из первого запроса

	hideAll := func() {
		r1.Hide()
		//r2.Hide()
		//r3.Hide()
		//r4.Hide()
		//r5.Hide()
		//r6.Hide()
		//r7.Hide()
		//r8.Hide()
		//r9.Hide()
		//r10.Hide()
	}
	hideAll()
	r1.Show()

	rs := []string{
		"Фильмы по цене",
		//"Врачи лечившие конкретного пациента",
		//"Пациенты лечившиеся у конкретного врача",
		//"Пациенты с одного участка",
		//"Пациенты, пользующиеся одной страховой компанией",
		//"Пациенты за определенный период по дате их рождения",
		//"Пациенты, лечившиеся у докторов определенной специализации",
		//"Пациенты, лечившиеся у докторов из одного отделения",
		//"Пациенты одного пола",
		//"Пациенты, живущие по одному и тому же месту жительствая",
	}
	requests := widget.NewSelect(rs, func(r string) {
		hideAll()
		switch r {
		case rs[0]:
			r1.Show()
			filmS.SetSelectedIndex(0)
			//case rs[1]:
			//	r2.Show()
			//	patS.SetSelectedIndex(0)
			//case rs[2]:
			//	r3.Show()
			//	docS.SetSelectedIndex(0)
			//case rs[3]:
			//	r4.Show()
			//	depS.SetSelectedIndex(0)
			//case rs[4]:
			//	r5.Show()
			//	cmpS.SetSelectedIndex(0)
			//case rs[5]:
			//	r6.Show()
			//	brtS.SetSelectedIndex(0)
			//case rs[6]:
			//	r7.Show()
			//	spcS.SetSelectedIndex(0)
			//case rs[7]:
			//	r8.Show()
			//	otdS.SetSelectedIndex(0)
			//case rs[8]:
			//	r9.Show()
			//	sexS.SetSelectedIndex(0)
			//case rs[9]:
			//	r10.Show()
			//	adrS.SetSelectedIndex(0)
		}
	})
	requests.SetSelectedIndex(0)

	requestsC := container.NewVBox(
		requests,
		r1,
		//r2,
		//r3,
		//r4,
		//r5,
		//r6,
		//r7,
		//r8,
		//r9,
		//r10,
		requestE,
	)
	requestsC.Move(fyne.NewPos(0.0, 300.0))
	requestsC.Resize(fyne.NewSize(1600, 200))
	requestsC.Refresh()

	mylog.Write([]byte("Request fields ready\n"))

	tabs := container.NewAppTabs(
		container.NewTabItem("Пациенты", patientsT),
		container.NewTabItem("Запросы", container.NewWithoutLayout(
			container.NewHBox(
				widget.NewButton("Запрос", func() {
					go func() {
						result.SetText(<-ch)
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

	prog := "Программа по производственной практике №2 \"Разработка программы для работы с базой данных\""
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
