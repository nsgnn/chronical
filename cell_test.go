package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestCell(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	testCases := []struct {
		name        string
		cell        *Cell
		action      func(c *Cell)
		expectedLog string
	}{
		{
			name: "NewCell Empty",
			action: func(c *Cell) {
				NewCell(0, 0, nil, nil)
			},
			expectedLog: `event="new_cell" x=0 y=0 value=  state=1`,
		},
		{
			name: "NewCell Given",
			action: func(c *Cell) {
				val := '5'
				NewCell(1, 1, &val, nil)
			},
			expectedLog: `event="new_cell" x=1 y=1 value=5 state=0`,
		},
		{
			name: "EnterValue",
			cell: NewCell(0, 0, nil, nil),
			action: func(c *Cell) {
				c.EnterValue('7')
			},
			expectedLog: `event="enter_value_success" x=0 y=0 value=7`,
		},
		{
			name: "EnterValue on Given Cell",
			cell: func() *Cell {
				val := '3'
				return NewCell(0, 0, &val, nil)
			}(),
			action: func(c *Cell) {
				c.EnterValue('7')
			},
			expectedLog: `event="enter_value_failed" reason="cell is given"`,
		},
		{
			name: "Clear",
			cell: func() *Cell {
				c := NewCell(0, 0, nil, nil)
				c.EnterValue('9')
				return c
			}(),
			action: func(c *Cell) {
				c.Clear()
			},
			expectedLog: `event="clear_success" x=0 y=0 value=9`,
		},
		{
			name: "Clear on Given Cell",
			cell: func() *Cell {
				val := '3'
				return NewCell(0, 0, &val, nil)
			}(),
			action: func(c *Cell) {
				c.Clear()
			},
			expectedLog: `event="clear_failed" reason="cell is given"`,
		},
		{
			name: "Validation Fails",
			cell: func() *Cell {
				c := NewCell(0, 0, nil, nil)
				c.EnterValue('4')
				return c
			}(),
			action: func(c *Cell) {
				c.RunValidation(false)
			},
			expectedLog: `event="validation_failed" x=0 y=0 value=4`,
		},
		{
			name: "Validation Passes",
			cell: func() *Cell {
				c := NewCell(0, 0, nil, nil)
				c.EnterValue('4')
				c.RunValidation(false)
				return c
			}(),
			action: func(c *Cell) {
				c.RunValidation(true)
			},
			expectedLog: `event="validation_passed" x=0 y=0 value=4`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.action(tc.cell)
			logOutput := strings.TrimSpace(buf.String())
			if !strings.Contains(logOutput, tc.expectedLog) {
				t.Errorf("expected log to contain %q, but got %q", tc.expectedLog, logOutput)
			}
		})
	}
}
