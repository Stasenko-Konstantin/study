package src

import (
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"reflect"
	"strconv"
	_ "strings"
)

var (
	patients = [][]string{[]string{"ID", "Страховая компания", "ФИО", "Место жительства",
		"Дата рождения", "Пол", "Участок"}}
	doctors    = [][]string{[]string{"ФИО", "Отделение", "Специализация", "Участок"}}
	talons = [][]string{[]string{"ID", "Дата и время приема", "Доктор", "Пациент"}}
)

func Start() {
	newLogger()

	defer func() {
		if r := recover(); r != nil {
			s := r.(error)
			log.Printf(s.Error() + " aga")
		}
	}()

	a := app.New()
	w := a.NewWindow("практика")
	w.Resize(fyne.NewSize(1800, 900))

	connect()
	log.Printf("DB connection complete")

	for _, e := range dbs {
		switch e.(type) {
		case []patient:
			for _, p := range e.([]patient) {
				var sex string
				if p.sex {
					sex = "Ж."
				} else {
					sex = "М."
				}
				patients = append(patients, []string{strconv.Itoa(p.id), p.insuranceCompany, p.sfm, p.residence,
					p.birth, sex, strconv.Itoa(p.district)})
			}
		case []doctor:
			for _, d := range e.([]doctor) {
				doctors = append(doctors, []string{d.sfm, d.department, d.specialization, strconv.Itoa(int(d.district.Int32))})
			}
		case []talon:
			for _, t := range e.([]talon) {
				talons = append(talons, []string{strconv.Itoa(t.id), t.reception, t.doctor, t.patient})
			}
		}
	}

	log.Printf("Data is full")

	patientsT := widget.NewTable(
		func() (int, int) {
			return len(patients), len(patients[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				o.(*fyne.Container).Objects[0] = widget.NewLabel(patients[i.Row][i.Col])
			} else {
				o.(*fyne.Container).Objects[0].(*widget.Entry).SetText(patients[i.Row][i.Col])
			}
		})
	patientsT.SetColumnWidth(0, 50.0)
	patientsT.SetColumnWidth(1, 170.0)
	patientsT.SetColumnWidth(2, 300.0)
	patientsT.SetColumnWidth(3, 430.0)
	patientsT.SetColumnWidth(4, 190.0)
	patientsT.SetColumnWidth(5, 50.0)
	patientsT.SetColumnWidth(6, 70.0)

	doctorsT := widget.NewTable(
		func() (int, int) {
			return len(doctors), len(doctors[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewEntry())
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				o.(*fyne.Container).Objects[0] = widget.NewLabel(doctors[i.Row][i.Col])
			} else {
				o.(*fyne.Container).Objects[0].(*widget.Entry).SetText(doctors[i.Row][i.Col])
			}
		})
	doctorsT.SetColumnWidth(0, 300.0)
	doctorsT.SetColumnWidth(1, 270.0)
	doctorsT.SetColumnWidth(2, 200.0)
	doctorsT.SetColumnWidth(3, 50.0)

	talonsT := widget.NewTable(
		func() (int, int) {
			return len(talons), len(talons[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(talons[i.Row][i.Col])
		})
	talonsT.SetColumnWidth(0, 50.0)
	talonsT.SetColumnWidth(1, 200.0)
	talonsT.SetColumnWidth(2, 300.0)
	talonsT.SetColumnWidth(3, 300.0)

	log.Printf("Table template ready")

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

	docs := makeDocs(doctors)
	docS := widget.NewSelect(docs, func(doc string) {
		requestE.SetText("P.sfm, P.birth, D.sfm from patients as P join doctors as D join talons as T on P.sfm = T.patient and D.sfm = T.doctor where P.district == " +
			findDoctor(doc, doctors))
	})
	docS.SetSelectedIndex(0)
	t1 := container.NewHBox(docS)
	t1Size := t1.Size()
	t1.Resize(fyne.NewSize(300, t1Size.Height))
	r1 := container.NewWithoutLayout(t1)
	docS.Move(fyne.NewPos(0, 10))
	docS.Refresh()

	r1.Show()
	docS.SetSelectedIndex(0)

	requestsC := container.NewVBox(
		r1,
		requestE,
	)
	requestsC.Move(fyne.NewPos(0.0, 300.0))
	requestsC.Resize(fyne.NewSize(1600, 200))
	requestsC.Refresh()

	log.Printf("Request fields ready")

	tabs := container.NewAppTabs(
		container.NewTabItem("Пациенты", patientsT),
		container.NewTabItem("Врачи", doctorsT),
		container.NewTabItem("Талоны", talonsT),
		container.NewTabItem("Запросы", container.NewWithoutLayout(
			container.NewHBox(
				widget.NewButton("Запрос", func() {
					go func() {
						defer func() {
							if r := recover(); r != nil {
								s := reflect.ValueOf(r)
								log.Printf(s.String())
								dialog.ShowError(errors.New(s.String()+"\nскорее всего неверный синтаксис ¯\\_:-/_/¯"), w)
							}
						}()
						db, err := sql.Open("sqlite3", "db.db")
						checkErr(err)
						rows, err := db.Query("select " + requestE.Text)
						checkErr(err)
						var r []*[]*sql.NullString
						i := 0
						for rows.Next() {
							var t []*sql.NullString
							r = append(r, &t)
							err := rows.ScanString(r[i])
							checkErr(err)
							i += 1
						}
						res := ""
						for _, e := range r {
							for _, ee := range *e {
								res += (*ee).String + " "
							}
							res += "\n"
						}
						result.SetText(res)
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

	log.Printf("Start render")

	w.SetContent(tabs)
	w.ShowAndRun()
}
