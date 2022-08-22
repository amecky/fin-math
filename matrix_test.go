package main

import (
	"testing"

	"github.com/amecky/fin-math/assert"
	"github.com/amecky/fin-math/math"
)

func TestMatrix(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(0, 1.0)
	m.AddRow("2").Set(0, 2.0)
	if m.Rows != 2 {
		t.Fatalf("Wrong number of rows - expected 2 but got %d", m.Rows)
	}
}

func TestMatrixDuplicateKey(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(0, 1.0)
	m.AddRow("1").Set(0, 2.0)
	if m.Rows != 1 {
		t.Fatalf("Wrong number of rows - expected 1 but got %d", m.Rows)
	}
	r := m.FindRow("1")
	if r == nil {
		t.Fatal("Expected row not found")
	}
	if r.Get(0) != 2.0 {
		t.Fatalf("Wrong value - expected 2.0 but got %.2f", r.Get(0))
	}
}

func TestMatrixIndexOutofBounds(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(10, 1.0)
	r := m.FindRow("1")
	if r == nil {
		t.Fatal("Expected row not found")
	}
	if r.Get(0) != 0.0 {
		t.Fatalf("Wrong value - expected 2.0 but got %.2f", r.Get(0))
	}
}

func TestMatrixIndexUnknownKey(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(10, 1.0)
	r := m.FindRow("xxxx")
	if r != nil {
		t.Fatal("No row should be found")
	}
}

func TestMatrixIndexMinMax(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(0, 1.0)
	m.AddRow("2").Set(0, 3.0)
	m.AddRow("3").Set(0, 8.0)
	m.AddRow("4").Set(0, 11.0)
	m.AddRow("5").Set(0, 5.0)
	m.AddRow("6").Set(0, 7.0)
	min, max := m.FindMinMaxBetween(0, 0, 20)
	assert.Equals(t, "Wrong min", 1.0, min)
	assert.Equals(t, "Wrong max", 11.0, max)
	min, max = m.FindMinMaxBetween(0, 3, 20)
	assert.Equals(t, "Wrong min", 5.0, min)
	assert.Equals(t, "Wrong max", 11.0, max)
	min, max = m.FindMinMaxBetween(0, 13, 20)
	assert.Equals(t, "Wrong min", 7.0, min)
	assert.Equals(t, "Wrong max", 7.0, max)
}

func TestMatrixIndexSum(t *testing.T) {
	m := math.NewMatrix(3)
	m.AddRow("1").Set(0, 1.0).Set(1, 3.0).Set(2, 8.0)
	r := m.FindRow("1")
	assert.NotNull(t, "Expected row not found", r)
	assert.Equals(t, "Wrong sum", 12.0, r.Sum())
}

func TestMatrixLast(t *testing.T) {
	m := math.NewMatrix(3)
	m.AddRow("1").Set(0, 1.0).Set(1, 3.0).Set(2, 8.0)
	r := m.Last()
	assert.NotNull(t, "Expected row not found", r)
	assert.Equals(t, "Wrong sum", 12.0, r.Sum())
}

func TestEmptyMatrixLast(t *testing.T) {
	m := math.NewMatrix(3)
	r := m.Last()
	if r != nil {
		t.Fatal("Row found")
	}
}

func TestMatrixAddColumn(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(0, 1.0)
	col := m.AddColumn()
	assert.Equals(t, "Wrong index", 1.0, float64(col))
	m.FindRow("1").Set(col, 2.0)
	row := m.FindRow("1")
	assert.Equals(t, "Wrong value", 2.0, row.Get(col))
}

func TestMatrixAdd(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(0, 1.0)
	m.AddRow("2").Set(0, 2.0)
	other := math.NewMatrix(1)
	other.AddRow("1").Set(0, 5.0)
	other.AddRow("2").Set(0, 6.0)
	col := m.Add(0, other, 0)
	assert.Equals(t, "Wrong column index", 1.0, float64(col))
	assert.Equals(t, "Wrong value", 8.0, m.FindRow("2").Get(1))
}
