package connect4

import (
	"fmt"
)

type Cell struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

type cellFunc func(Cell) Cell

func (c Cell) Left() Cell {
	return Cell{c.Col - 1, c.Row}
}

func (c Cell) Right() Cell {
	return Cell{c.Col + 1, c.Row}
}

func (c Cell) Top() Cell {
	return Cell{c.Col, c.Row + 1}
}

func (c Cell) Bottom() Cell {
	return Cell{c.Col, c.Row - 1}
}

func (c Cell) TopLeft() Cell {
	return c.Top().Left()
}

func (c Cell) TopRight() Cell {
	return c.Top().Right()
}

func (c Cell) BotLeft() Cell {
	return c.Bottom().Left()
}

func (c Cell) BotRight() Cell {
	return c.Bottom().Right()
}

type Field [][]Color

func NewField(cols, rows int) Field {
	field := make(Field, rows)
	for i := range field {
		field[i] = make([]Color, cols)
	}

	return field
}

func (f Field) InRange(c Cell) bool {
	return (c.Row >= 0 && c.Row < len(f)) &&
		(c.Col >= 0 && c.Col < len(f[0]))
}

func (f Field) Equal(c1, c2 Cell) bool {
	if f.InRange(c1) && f.InRange(c2) {
		return f[c1.Row][c1.Col] == f[c2.Row][c2.Col]
	}

	// out of range
	return false
}

func (f Field) Full() bool {
	for _, cell := range f[0] {
		if cell == NullColor {
			return false
		}
	}
	return true
}

func (f Field) Empty() bool {
	lastRow := len(f) - 1

	for _, cell := range f[lastRow] {
		if cell != NullColor {
			return false
		}
	}
	return true
}

func (f *Field) Update(col int, color Color) (Cell, error) {
	if !f.InRange(Cell{Col: col}) {
		return Cell{}, fmt.Errorf("col %d is out of range", col)
	}

	for row := len(*f) - 1; row >= 0; row-- {
		if (*f)[row][col] == NullColor {
			(*f)[row][col] = color
			return Cell{col, row}, nil
		}
	}

	return Cell{}, fmt.Errorf("col %d is full", col)
}

func (f *Field) FindFour(c Cell) []Cell {
	if seq := f.findHelper(c, Cell.Left, Cell.Right); len(seq) >= 4 {
		return seq
	}
	if seq := f.findHelper(c, Cell.Top, Cell.Bottom); len(seq) >= 4 {
		return seq
	}
	if seq := f.findHelper(c, Cell.TopLeft, Cell.BotRight); len(seq) >= 4 {
		return seq
	}
	if seq := f.findHelper(c, Cell.TopRight, Cell.BotLeft); len(seq) >= 4 {
		return seq
	}

	return nil
}

func (f *Field) findHelper(c Cell, f1, f2 cellFunc) []Cell {
	sequence := []Cell{c}

	for cI := f1(c); f.Equal(c, cI); cI = f1(cI) {
		sequence = append(sequence, cI)
	}
	for cI := f2(c); f.Equal(c, cI); cI = f2(cI) {
		sequence = append(sequence, cI)
	}

	return sequence
}
