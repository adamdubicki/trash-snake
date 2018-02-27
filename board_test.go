package main

import "testing"

func TestBoardConstructor(t *testing.T) {
	b := createBoard(3, 4)
	if b.Width != 3 {
		t.Errorf("Board width was incorrect, got: %d, want: %d.", b.Width, 3)
	}

	if b.Height != 4 {
		t.Errorf("Board height was incorrect, got: %d, want: %d.", b.Height, 4)
	}
}

func TestXInBounds(t *testing.T) {
	b := createBoard(3, 4)
	tables := []struct {
		x        int
		inBounds bool
	}{
		{-1, false},
		{0, true},
		{1, true},
		{2, true},
		{3, faelse},
	}
	for _, table := range tables {
		inBounds := b.xInBounds(table.x)
		if inBounds != table.inBounds {
			t.Errorf("xInbounds of (%d) was incorrect, got: %t, want: %t.", table.x, inBounds, table.inBounds)
		}
	}
}

func TestYInBounds(t *testing.T) {
	b := createBoard(3, 4)
	tables := []struct {
		y        int
		inBounds bool
	}{
		{-1, false},
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, false},
	}
	for _, table := range tables {
		inBounds := b.yInBounds(table.y)
		if inBounds != table.inBounds {
			t.Errorf("yInbounds of (%d) was incorrect, got: %t, want: %t.", table.y, inBounds, table.inBounds)
		}
	}
}

func TestGetTile(t *testing.T) {
	b := createBoard(3, 4)
	b.insert(Point{1, 2}, obstacle())
	tables := []struct {
		p Point
		e int
	}{
		{Point{0, 0}, empty().EntityType},
		{Point{-1, -1}, invalid().EntityType},
		{Point{1, 2}, obstacle().EntityType},
	}
	for _, table := range tables {
		tileEntity := b.getTile(table.p)
		if tileEntity.EntityType != table.e {
			t.Errorf("Get tile failed, Got (%d) wanted (%d)", tileEntity.EntityType, table.e)
		}
	}
}
