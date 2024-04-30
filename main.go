package main

import (
	"bufio"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	uniqueid "github.com/albinj12/unique-id"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Student struct {
	id         int
	name       string
	surname    string
	patronymic string
}

type Variant struct {
	id         int
	pathToFile string
}

type TableRow struct {
	studentId int
	variantId int
	mark      int
}

var students = make(map[int]Student)
var variants = make(map[int]Variant)
var table = make(map[int]TableRow)

var a = app.New()
var window = a.NewWindow("Database")
var buttons = container.NewWithoutLayout()

func main() {
	a.Settings().SetTheme(theme.DarkTheme())
	window.Resize(fyne.NewSize(1920, 1080))
	//window.SetFullScreen(true)
	ic, _ := fyne.LoadResourceFromPath("Sisyphus.jpg")
	window.SetIcon(ic)

	class := ""

	temp := container.NewVBox()
	createDB := widget.NewButton("Create new database", func() {
		input := widget.NewEntry()
		input.SetPlaceHolder("Enter new database name:")
		temp.Add(input)
		temp.Add(widget.NewButton("Confirm:", func() {
			class = input.Text
			os.Mkdir(class, os.ModePerm)
			os.Create(fmt.Sprintf("%s/names.txt", class))
			os.Create(fmt.Sprintf("%s/variants.txt", class))
			file, _ := os.OpenFile(fmt.Sprintf("%s/variants.txt", class), os.O_RDWR, 0644)
			file.WriteString("1 var1\n2 var2\n3 var3\n4 var4\n5 var5\n")
			os.Create(fmt.Sprintf("%s/marks.txt", class))
			os.Create(fmt.Sprintf("%s/backup.txt", class))
			setup(class)
		}))
	})
	temp.Add(createDB)

	openDB := widget.NewButton("Open existing database", func() {
		input := widget.NewEntry()
		input.SetPlaceHolder("Enter database name:")
		temp.Add(input)
		temp.Add(widget.NewButton("Confirm:", func() {
			class = input.Text
			setup(class)
		}))
	})

	temp.Add(openDB)
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), temp))
	window.ShowAndRun()
}

func setup(classNumber string) {
	names, _ := os.Open(fmt.Sprintf("%s/names.txt", classNumber))
	defer names.Close()

	scanner := bufio.NewScanner(names)
	for scanner.Scan() {
		line := scanner.Text()
		params := strings.Fields(line)

		id, _ := strconv.Atoi(params[0])
		students[id] = Student{id,
			params[1],
			params[2],
			params[3]}
	}

	vars, _ := os.Open(fmt.Sprintf("%s/variants.txt", classNumber))
	defer vars.Close()

	scanner = bufio.NewScanner(vars)
	for scanner.Scan() {
		line := scanner.Text()
		params := strings.Fields(line)

		id, _ := strconv.Atoi(params[0])
		variants[id] = Variant{id, params[1]}
	}

	file, _ := os.OpenFile(fmt.Sprintf("%s/marks.txt", classNumber), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	marks, _ := os.Open(fmt.Sprintf("%s/marks.txt", classNumber))
	fi, _ := os.Stat(fmt.Sprintf("%s/marks.txt", classNumber))

	if fi.Size() == 0 {
		for _, s := range students {
			i := rand.Intn(len(variants))
			for _, v := range variants {
				if i == 0 {
					table[s.id] = TableRow{s.id, v.id, 0}
					file.WriteString(fmt.Sprintln(s.id, v.id, 0))
				}
				i--
			}
		}
	} else {
		scanner := bufio.NewScanner(marks)

		for scanner.Scan() {
			line := scanner.Text()
			params := strings.Fields(line)
			studentId, _ := strconv.Atoi(params[0])
			variantId, _ := strconv.Atoi(params[1])
			mark, _ := strconv.Atoi(params[2])
			table[studentId] = TableRow{studentId, variantId, mark}
		}
	}

	fyneSetup(classNumber)
}

func fyneSetup(class string) {
	show := widget.NewButton("Show", func() {
		showAll()
	})
	student := widget.NewButton("Show student", func() {
		showStudent()
	})
	add := widget.NewButton("Add student", func() {
		addStudent(class)
	})
	delete := widget.NewButton("Delete student", func() {
		deleteStudent(class)
	})
	update := widget.NewButton("Update mark", func() {
		update(class)
	})

	load := widget.NewButton("Load backup", func() {
		loadBackup(class)
	})

	backup := widget.NewButton("Make backup", func() {
		makeBackup(class)
	})

	close := widget.NewButton("Close program", func() {
		closeProgram()
	})

	buttons = container.NewVBox(show, student, add, delete, update, load, backup, close)

	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), buttons))
}

func getTempWithInput() (*widget.Entry, fyne.Container) {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter full name:")
	temp := *buttons
	temp.Add(input)

	return input, temp
}

func showAll() {
	strAll := ""
	for _, v := range table {
		strAll += fmt.Sprintf("%s: %s, %d\n",
			fmt.Sprintf("%s %s %s", students[v.studentId].name,
				students[v.studentId].surname, students[v.studentId].patronymic),
			variants[v.variantId].pathToFile, v.mark)
	}
	grid := widget.NewTextGridFromString(strAll)
	temp := *buttons
	scrollableGrid := container.NewScroll(grid)
	scrollableGrid.SetMinSize(fyne.NewSize(200, 400))
	temp.Add(scrollableGrid)
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func showStudent() {
	studentLine := widget.NewLabel("")
	input, temp := getTempWithInput()
	temp.Add(widget.NewButton("Find", func() {
		temp.Remove(studentLine)
		params := strings.Fields(input.Text)

		if len(params) != 3 {
			studentLine = widget.NewLabel("Wrong number of parameters.")
			temp.Add(studentLine)
			return
		} else {
			name, surname, patronymic := params[0], params[1], params[2]

			for _, v := range table {
				if students[v.studentId].name == name &&
					students[v.studentId].surname == surname &&
					students[v.studentId].patronymic == patronymic {
					studentLine = widget.NewLabel(fmt.Sprintf("%s: %s, %d",
						fmt.Sprintf("%s %s %s", students[v.studentId].name,
							students[v.studentId].surname, students[v.studentId].patronymic),
						variants[v.variantId].pathToFile, v.mark))
					temp.Add(studentLine)
					return
				}
			}
			studentLine = widget.NewLabel("Didn't find student with this name.")
		}
		temp.Add(studentLine)
	}))
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func addStudent(class string) {
	file, _ := os.OpenFile(fmt.Sprintf("%s/names.txt", class), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	//defer file.Close()

	studentLine := widget.NewLabel("")
	input, temp := getTempWithInput()

	temp.Add(widget.NewButton("Add", func() {
		temp.Remove(studentLine)
		params := strings.Fields(input.Text)

		if len(params) != 3 {
			studentLine = widget.NewLabel("Wrong number of parameters.")
			temp.Add(studentLine)
			return
		} else {
			name, surname, patronymic := params[0], params[1], params[2]

			for _, v := range students {
				if v.name == name && v.surname == surname && v.patronymic == patronymic {
					studentLine = widget.NewLabel("Student with this name already exists.")
					temp.Add(studentLine)
					return
				}
			}

			id, _ := uniqueid.Generateid("n", 9)
			newId, _ := strconv.Atoi(id)
			students[newId] = Student{newId, name, surname, patronymic}
			file.WriteString(fmt.Sprintln(newId, name, surname, patronymic))
			i := rand.Intn(len(variants))
			for _, v := range variants {
				if i == 0 {
					marks, _ := os.OpenFile(fmt.Sprintf("%d/marks.txt", class), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
					table[students[newId].id] = TableRow{students[newId].id, v.id, 0}
					marks.WriteString(fmt.Sprintln(students[newId].id, v.id, 0))
				}
				i--
			}
			studentLine = widget.NewLabel("Student added successfully!")
		}
		temp.Add(studentLine)
	}))
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func deleteStudent(class string) {
	studentLine := widget.NewLabel("")
	input, temp := getTempWithInput()

	temp.Add(widget.NewButton("Delete", func() {
		temp.Remove(studentLine)
		params := strings.Fields(input.Text)

		if len(params) != 3 {
			studentLine = widget.NewLabel("Wrong number of parameters.")
			temp.Add(studentLine)
			return
		} else {
			name, surname, patronymic := params[0], params[1], params[2]

			for _, v := range students {
				if v.name == name && v.surname == surname && v.patronymic == patronymic {
					studentLine = widget.NewLabel("Student removed successfully.")
					delete(students, v.id)
					delete(table, v.id)

					file, _ := os.OpenFile(fmt.Sprintf("%s/names.txt", class), os.O_RDWR, 0644)
					scanner := bufio.NewScanner(file)
					var bs1 []byte
					buf := bytes.NewBuffer(bs1)

					var text string
					for scanner.Scan() {
						text = scanner.Text()
						if !(strings.Contains(text, fmt.Sprintf("%s %s %s", name, surname, patronymic))) {
							buf.WriteString(text + "\n")
						}
					}
					file.Truncate(0)
					file.Seek(0, 0)
					buf.WriteTo(file)

					marks, _ := os.OpenFile(fmt.Sprintf("%s/marks.txt", class), os.O_RDWR, 0644)
					scanner = bufio.NewScanner(marks)
					var bs2 []byte
					buf = bytes.NewBuffer(bs2)
					for scanner.Scan() {
						text = scanner.Text()
						if !(strings.Contains(text, fmt.Sprintf("%d", v.id))) {
							buf.WriteString(text + "\n")
						}
					}
					marks.Truncate(0)
					marks.Seek(0, 0)
					buf.WriteTo(marks)

					temp.Add(studentLine)
					return
				}
			}

			studentLine = widget.NewLabel("Didn't find student with this name.")
		}
		temp.Add(studentLine)
	}))
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func update(class string) {
	studentLine := widget.NewLabel("")
	input, temp := getTempWithInput()

	temp.Add(widget.NewButton("Enter", func() {
		temp.Remove(studentLine)
		params := strings.Fields(input.Text)

		if len(params) != 3 {
			studentLine = widget.NewLabel("Wrong number of parameters.")
			temp.Add(studentLine)
			return
		} else {
			name, surname, patronymic := params[0], params[1], params[2]

			file, _ := os.OpenFile(fmt.Sprintf("%s/marks.txt", class), os.O_RDWR, 0644)

			for _, v := range students {
				if v.name == name && v.surname == surname && v.patronymic == patronymic {
					inputMark := widget.NewEntry()
					inputMark.SetPlaceHolder("Enter mark:")
					temp.Add(inputMark)
					temp.Add(widget.NewButton("Update", func() {
						temp.Remove(studentLine)
						mark, _ := strconv.Atoi(inputMark.Text)
						changedRow := table[v.id]
						changedRow.mark = mark
						table[v.id] = changedRow

						scanner := bufio.NewScanner(file)
						var bs []byte
						buf := bytes.NewBuffer(bs)

						var text string
						for scanner.Scan() {
							text = scanner.Text()
							//fmt.Println(text, "  :  ", fmt.Sprintf("%d %d %d", table[v.id].studentId, table[v.id].variantId, prevMark))
							if strings.Contains(text, fmt.Sprintf("%d", v.id)) {
								buf.WriteString(fmt.Sprintf("%d %d %d\n", v.id, table[v.id].variantId, mark))
							} else {
								buf.WriteString(text + "\n")
							}
						}
						file.Truncate(0)
						file.Seek(0, 0)
						buf.WriteTo(file)
						studentLine = widget.NewLabel("Mark changed successfully!")
						temp.Add(studentLine)
						return
					}))
					temp.Add(studentLine)
					return
				}
			}
			studentLine = widget.NewLabel("Didn't find student with this name.")
		}
		temp.Add(studentLine)
	}))
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func loadBackup(class string) {
	file, _ := os.Open(fmt.Sprintf("%s/backup.txt", class))

	studentLine := widget.NewLabel("")
	temp := *buttons

	scanner := bufio.NewScanner(file)

	names, _ := os.OpenFile(fmt.Sprintf("%s/names.txt", class), os.O_RDWR, 0644)
	vars, _ := os.OpenFile(fmt.Sprintf("%s/variants.txt", class), os.O_RDWR, 0644)
	marks, _ := os.OpenFile(fmt.Sprintf("%s/marks.txt", class), os.O_RDWR, 0644)
	names.Truncate(0)
	vars.Truncate(0)
	marks.Truncate(0)

	for k := range students {
		delete(students, k)
	}
	for k := range variants {
		delete(variants, k)
	}
	for k := range table {
		delete(table, k)
	}

	//writer := bufio.NewWriter(names)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		names.WriteString(line + "\n")
		params := strings.Fields(line)

		id, _ := strconv.Atoi(params[0])
		students[id] = Student{id,
			params[1],
			params[2],
			params[3]}
	}

	//writer = bufio.NewWriter(vars)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		vars.WriteString(line + "\n")

		params := strings.Fields(line)
		id, _ := strconv.Atoi(params[0])
		variants[id] = Variant{id, params[1]}
	}

	//writer = bufio.NewWriter(marks)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		marks.WriteString(line + "\n")

		params := strings.Fields(line)
		studentId, _ := strconv.Atoi(params[0])
		variantId, _ := strconv.Atoi(params[1])
		mark, _ := strconv.Atoi(params[2])
		table[studentId] = TableRow{studentId, variantId, mark}
	}

	studentLine = widget.NewLabel("Backup loaded successfully!")
	temp.Add(studentLine)
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func makeBackup(class string) {
	file, _ := os.OpenFile(fmt.Sprintf("%s/backup.txt", class), os.O_RDWR, 0644)

	studentLine := widget.NewLabel("")
	temp := *buttons

	file.Truncate(0)

	writer := bufio.NewWriter(file)

	for _, v := range students {
		writer.WriteString(fmt.Sprintf("%d %s %s %s\n", v.id, v.name, v.surname, v.patronymic))
	}

	writer.WriteString("\n")

	for _, v := range variants {
		writer.WriteString(fmt.Sprintf("%d %s\n", v.id, v.pathToFile))
	}

	writer.WriteString("\n")

	for _, v := range table {
		writer.WriteString(fmt.Sprintf("%d %d %d\n", v.studentId, v.variantId, v.mark))
	}
	writer.Flush()

	studentLine = widget.NewLabel("Backup created successfully!")
	temp.Add(studentLine)
	window.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewImageFromFile("Sisyphus.jpg"), &temp))
}

func closeProgram() {
	window.Close()
}
