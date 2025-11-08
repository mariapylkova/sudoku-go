package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	rows     = 9
	columns  = 9
	cellSize = 50
)

type Grid struct {
	size  [rows][columns]int8
	fixed [rows][columns]bool
	cells [rows][columns]fyne.CanvasObject
}

func NewSudoku(digits [rows][columns]int8) *Grid {
	var grid Grid
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			grid.size[r][c] = digits[r][c]
			if digits[r][c] != 0 {
				grid.fixed[r][c] = true
			}
		}
	}
	return &grid
}

func (g *Grid) ValidNumberPosition(number int8, row, column int) error {
	if g.fixed[row][column] {
		return fmt.Errorf("ячейка фиксирована")
	}
	for i := 0; i < 9; i++ {
		if g.size[row][i] == number {
			return fmt.Errorf("число уже есть в строке")
		}
		if g.size[i][column] == number {
			return fmt.Errorf("число уже есть в столбце")
		}
	}
	startRow := (row / 3) * 3
	startCol := (column / 3) * 3
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			if g.size[i][j] == number {
				return fmt.Errorf("число уже есть в квадрате 3x3")
			}
		}
	}
	return nil
}

func (g *Grid) WriteNumber(number int8, row, column int) error {
	if number < 1 || number > 9 {
		return fmt.Errorf("число должно быть от 1 до 9")
	}
	if err := g.ValidNumberPosition(number, row, column); err != nil {
		return err
	}
	g.size[row][column] = number

	obj := g.cells[row][column]
	if obj == nil {
		return nil
	}
	switch w := obj.(type) {
	case *widget.Button:
		w.SetText(strconv.Itoa(int(number)))
	case *widget.Label:
		w.SetText(strconv.Itoa(int(number)))
	default:
	}
	return nil
}

func (g *Grid) IsComplete() bool {
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			if g.size[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

func (g *Grid) SetupUI(win fyne.Window) fyne.CanvasObject {
	canvasContainer := container.NewWithoutLayout()

	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			val := g.size[r][c]
			display := ""
			if val != 0 {
				display = fmt.Sprintf("%d", val)
			}

			row, col := r, c

			var cellObj fyne.CanvasObject

			if g.fixed[r][c] {
				lbl := widget.NewLabelWithStyle(display, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
				lbl.Resize(fyne.NewSize(cellSize, cellSize))
				lbl.Move(fyne.NewPos(float32(c*cellSize), float32(r*cellSize)))
				cellObj = lbl
			} else {
				var btn *widget.Button
				var actionDialog dialog.Dialog

				btn = widget.NewButton(display, func() {

					actionDialog = dialog.NewCustom("Выберите действие", "Закрыть",
						container.NewVBox(
							widget.NewButton("Записать число", func() {

								entry := widget.NewEntry()
								entry.SetPlaceHolder("Введите число 1–9")

								dialog.NewForm(
									"Введите число",
									"OK",
									"Отмена",
									[]*widget.FormItem{widget.NewFormItem("", entry)},
									func(ok bool) {
										if !ok {
											return
										}
										n, err := strconv.Atoi(entry.Text)
										if err != nil {
											dialog.ShowError(fmt.Errorf("неверный ввод"), win)
											return
										}
										if err := g.WriteNumber(int8(n), row, col); err != nil {
											dialog.ShowError(err, win)
										} else {
											btn.SetText(entry.Text)
											actionDialog.Hide()
										}
										if g.IsComplete() {
											dialog.ShowInformation("Победа!", "Вы успешно решили судоку 🎉", win)
										}
									},
									win,
								).Show()
							}),

							widget.NewButton("Очистить", func() {
								g.size[row][col] = 0
								btn.SetText("")
								actionDialog.Hide()
							}),
						),
						win)

					actionDialog.Show()
				})

				btn.Resize(fyne.NewSize(cellSize, cellSize))
				btn.Move(fyne.NewPos(float32(c*cellSize), float32(r*cellSize)))
				cellObj = btn
			}

			g.cells[r][c] = cellObj
			canvasContainer.Add(cellObj)
		}
	}

	for i := 0; i <= 9; i++ {
		lineV := canvas.NewLine(color.White)
		lineH := canvas.NewLine(color.White)

		if i%3 == 0 {
			lineV.StrokeWidth = 3
			lineH.StrokeWidth = 3
		} else {
			lineV.StrokeWidth = 1
			lineH.StrokeWidth = 1
		}

		lineV.Move(fyne.NewPos(float32(i*cellSize), 0))
		lineV.Resize(fyne.NewSize(1, cellSize*rows))
		lineH.Move(fyne.NewPos(0, float32(i*cellSize)))
		lineH.Resize(fyne.NewSize(cellSize*columns, 1))

		canvasContainer.Add(lineV)
		canvasContainer.Add(lineH)
	}

	return canvasContainer
}

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Sudoku")

	startDigits := [rows][columns]int8{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	grid := NewSudoku(startDigits)
	content := grid.SetupUI(win)

	win.SetContent(content)
	win.Resize(fyne.NewSize(cellSize*columns+20, cellSize*rows+40))
	win.ShowAndRun()
}
