package main

import (
	"testing"
)

func TestNewCell(t *testing.T) {
	t.Run("given", func(t *testing.T) {
		g := '1'
		c := NewCell(0, 0, &g)
		if c.state != given {
			t.Errorf("expected state %d, got %d", given, c.state)
		}
		if c.value != g {
			t.Errorf("expected value %c, got %c", g, c.value)
		}
	})

	t.Run("empty", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		if c.state != empty {
			t.Errorf("expected state %d, got %d", empty, c.state)
		}
		if c.value != ' ' {
			t.Errorf("expected value ' ', got %c", c.value)
		}
	})
}

func TestEnterValue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		v := '1'
		c.EnterValue(v)
		if c.state != filled {
			t.Errorf("expected state %d, got %d", filled, c.state)
		}
		if c.value != v {
			t.Errorf("expected value %c, got %c", v, c.value)
		}
	})

	t.Run("given", func(t *testing.T) {
		g := '1'
		c := NewCell(0, 0, &g)
		v := '2'
		c.EnterValue(v)
		if c.state != given {
			t.Errorf("expected state %d, got %d", given, c.state)
		}
		if c.value != g {
			t.Errorf("expected value %c, got %c", g, c.value)
		}
	})
}

func TestClear(t *testing.T) {
	t.Run("filled", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		c.EnterValue('1')
		c.Clear()
		if c.state != empty {
			t.Errorf("expected state %d, got %d", empty, c.state)
		}
		if c.value != ' ' {
			t.Errorf("expected value ' ', got %c", c.value)
		}
	})

	t.Run("given", func(t *testing.T) {
		g := '1'
		c := NewCell(0, 0, &g)
		c.Clear()
		if c.state != given {
			t.Errorf("expected state %d, got %d", given, c.state)
		}
		if c.value != g {
			t.Errorf("expected value %c, got %c", g, c.value)
		}
	})

	t.Run("empty", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		c.Clear()
		if c.state != empty {
			t.Errorf("expected state %d, got %d", empty, c.state)
		}
		if c.value != ' ' {
			t.Errorf("expected value ' ', got %c", c.value)
		}
	})
}

func TestRunValidation(t *testing.T) {
	t.Run("filled_correct", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		c.EnterValue('1')
		c.RunValidation(true)
		if c.state != filled {
			t.Errorf("expected state %d, got %d", filled, c.state)
		}
	})

	t.Run("filled_incorrect", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		c.EnterValue('1')
		c.RunValidation(false)
		if c.state != invalid {
			t.Errorf("expected state %d, got %d", invalid, c.state)
		}
	})

	t.Run("invalid_correct", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		c.EnterValue('1')
		c.RunValidation(false)
		c.RunValidation(true)
		if c.state != filled {
			t.Errorf("expected state %d, got %d", filled, c.state)
		}
	})

	t.Run("invalid_incorrect", func(t *testing.T) {
		c := NewCell(0, 0, nil)
		c.EnterValue('1')
		c.RunValidation(false)
		c.RunValidation(false)
		if c.state != invalid {
			t.Errorf("expected state %d, got %d", invalid, c.state)
		}
	})
}
