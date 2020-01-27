package main

import "testing"

func TestTable(t *testing.T) {
	tab := NewTable()
	tab.Set(0, 5, "35")
	if got, want := tab.At(0, 5), "35"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	r := tab.Row(0)
	if got, want := r[5], "35"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	c := tab.FindColumn(0, "35")
	if got, want := c, 5; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	fr := tab.FindRow(5, "35")
	if got, want := fr, 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTable2(t *testing.T) {
	tab := NewTable()
	tab.Set(3, 5, "35")
	if got, want := tab.At(3, 5), "35"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	r := tab.Row(3)
	if got, want := r[5], "35"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	c := tab.FindColumn(3, "35")
	if got, want := c, 5; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	fr := tab.FindRow(5, "35")
	if got, want := fr, 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
