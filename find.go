package connect4

func (g *Game) findFour(c Cell) []Cell {
	if seq := g.findHelper(c, Cell.Left, Cell.Right); len(seq) == 4 {
		return seq
	}
	if seq := g.findHelper(c, Cell.Top, Cell.Bottom); len(seq) == 4 {
		return seq
	}
	if seq := g.findHelper(c, Cell.TopLeft, Cell.BotRight); len(seq) == 4 {
		return seq
	}
	if seq := g.findHelper(c, Cell.TopRight, Cell.BotLeft); len(seq) == 4 {
		return seq
	}

	return nil
}

func (g *Game) findHelper(c Cell, f1, f2 cellFunc) []Cell {
	sequence := []Cell{c}

	for cI := f1(c); g.field.Equal(c, cI); cI = f1(cI) {
		sequence = append(sequence, cI)
	}
	for cI := f2(c); g.field.Equal(c, cI); cI = f2(cI) {
		sequence = append(sequence, cI)
	}

	return sequence
}
