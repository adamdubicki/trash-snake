package main

import "testing"

func TestFood(t *testing.T) {
	tables := []struct {
		displayString string
		isDangerous   bool
		entityEnum    int
	}{
		{"f", false, FOOD},
	}
	for _, table := range tables {
		e := food()
		if e.Display != table.displayString {
			t.Errorf("Display string did not match (%d)", e)
		}
		if e.Dangerous != table.isDangerous {
			t.Errorf("Dangerous bool did not match (%d)", e)
		}
		if e.EntityType != table.entityEnum {
			t.Errorf("Type did not match (%d)", e)
		}
	}
}

func TestEmpty(t *testing.T) {
	tables := []struct {
		displayString string
		isDangerous   bool
		entityEnum    int
	}{
		{" ", false, EMPTY},
	}
	for _, table := range tables {
		e := empty()
		if e.Display != table.displayString {
			t.Errorf("Display string did not match (%d)", e)
		}
		if e.Dangerous != table.isDangerous {
			t.Errorf("Dangerous bool did not match (%d)", e)
		}
		if e.EntityType != table.entityEnum {
			t.Errorf("Type did not match (%d)", e)
		}
	}
}

func TestSnakeHead(t *testing.T) {
	tables := []struct {
		displayString string
		isDangerous   bool
		entityEnum    int
	}{
		{"h", true, SNAKEHEAD},
	}
	for _, table := range tables {
		e := snakeHead()
		if e.Display != table.displayString {
			t.Errorf("Display string did not match (%d)", e)
		}
		if e.Dangerous != table.isDangerous {
			t.Errorf("Dangerous bool did not match (%d)", e)
		}
		if e.EntityType != table.entityEnum {
			t.Errorf("Type did not match (%d)", e)
		}
	}
}

func TestObstacle(t *testing.T) {
	tables := []struct {
		displayString string
		isDangerous   bool
		entityEnum    int
	}{
		{"o", true, OBSTACLE},
	}
	for _, table := range tables {
		e := obstacle()
		if e.Display != table.displayString {
			t.Errorf("Display string did not match (%d)", e)
		}
		if e.Dangerous != table.isDangerous {
			t.Errorf("Dangerous bool did not match (%d)", e)
		}
		if e.EntityType != table.entityEnum {
			t.Errorf("Type did not match (%d)", e)
		}
	}
}

func TestInvalid(t *testing.T) {
	tables := []struct {
		displayString string
		isDangerous   bool
		entityEnum    int
	}{
		{"X", true, INVALID},
	}
	for _, table := range tables {
		e := invalid()
		if e.Display != table.displayString {
			t.Errorf("Display string did not match (%d)", e)
		}
		if e.Dangerous != table.isDangerous {
			t.Errorf("Dangerous bool did not match (%d)", e)
		}
		if e.EntityType != table.entityEnum {
			t.Errorf("Type did not match (%d)", e)
		}
	}
}
