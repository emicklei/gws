package main

// Table is a dynamic structure to hold strings indexed at row,column, zero-based.
type Table struct {
	data      map[uint]map[uint]string
	maxRow    int
	maxColumn int
}

func NewTable() *Table {
	return &Table{data: map[uint]map[uint]string{}}
}

func (t *Table) Rows() int {
	return t.maxRow + 1
}

func (t *Table) FindColumn(row int, value string) int {
	r, ok := t.data[uint(row)]
	if !ok {
		return -1
	}
	for column, each := range r {
		if each == value {
			return int(column)
		}
	}
	return -1
}

func (t *Table) FindRow(column int, value string) int {
	for r := 0; r <= t.maxRow; r++ {
		if c := t.FindColumn(r, value); c != -1 {
			return r
		}
	}
	return -1
}

func (t *Table) At(row, column int) string {
	r, ok := t.data[uint(row)]
	if !ok {
		return ""
	}
	v, ok := r[uint(column)]
	if !ok {
		return ""
	}
	return v
}

func (t *Table) Set(row, column int, value string) {
	r, ok := t.data[uint(row)]
	if !ok {
		r = map[uint]string{}
		t.data[uint(row)] = r
	}
	r[uint(column)] = value
	if row > t.maxRow {
		t.maxRow = row
	}
	if column > t.maxColumn {
		t.maxColumn = column
	}
}

func (t *Table) Row(row int) (list []string) {
	r, ok := t.data[uint(row)]
	if !ok {
		for s := 0; s <= t.maxColumn; s++ {
			list = append(list, "")
		}
		return
	}
	for c := 0; c <= t.maxColumn; c++ {
		v, ok := r[uint(c)]
		if !ok {
			list = append(list, "")
		} else {
			list = append(list, v)
		}
	}
	return
}
