package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNonogramEngine(t *testing.T) {
	level := Level{
		ID:       1,
		Name:     "Test Level",
		Author:   "Test Author",
		Initial:  "    \n 1  \n  1 \n    ",
		Solution: "1111\n1111\n1111\n1111",
		Engine:   "nonogram",
		Width:    4,
		Height:   4,
	}

	initialSave := &Save{
		LevelID: 1,
		State:   "    \n 1  \n  1 \n    ",
		Solved:  false,
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)

	testCases := []struct {
		name          string
		action        func(e *NonogramEngine)
		expectedState string
		expectedLog   string
		initialSave   *Save
	}{
		{
			name: "New Game",
			action: func(e *NonogramEngine) {
				// No action needed, just check initial state
			},
			expectedState: "    \n 1  \n  1 \n    ",
			expectedLog:   "event=\"StatefulLevelLoad\"",
			initialSave:   initialSave,
		},
		{
			name: "Primary Action",
			action: func(e *NonogramEngine) {
				e.PrimaryAction(0, 0)
			},
			expectedState: "1   \n 1  \n  1 \n    ",
			expectedLog:   "event=\"enter_value_success\" x=0 y=0 value=1",
			initialSave:   initialSave,
		},
		{
			name: "Secondary Action",
			action: func(e *NonogramEngine) {
				e.SecondaryAction(1, 0)
			},
			expectedState: " X  \n 1  \n  1 \n    ",
			expectedLog:   "event=\"enter_value_success\" x=0 y=1 value=X",
			initialSave:   initialSave,
		},
		{
			name: "Evaluate Incorrect Solution",
			action: func(e *NonogramEngine) {
				solved, err := e.EvaluateSolution()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if solved {
					t.Errorf("expected solution to be incorrect")
				}
			},
			expectedState: "    \n 1  \n  1 \n    ",
			expectedLog:   "",
			initialSave:   initialSave,
		},
		{
			name: "Evaluate Correct Solution",
			action: func(e *NonogramEngine) {
				e.Engine.Save.State = "1111\n1111\n1111\n1111"
				solved, err := e.EvaluateSolution()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !solved {
					t.Errorf("expected solution to be correct")
				}
			},
			expectedState: "1111\n1111\n1111\n1111",
			expectedLog:   "",
			initialSave: &Save{
				LevelID: 1,
				State:   "1111\n1111\n1111\n1111",
				Solved:  false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			engine := &NonogramEngine{}
			engine.New(level, tc.initialSave)

			if tc.action != nil {
				tc.action(engine)
			}

			if engine.Engine.Save.State != tc.expectedState {
				t.Errorf("expected state %q, got %q", tc.expectedState, engine.Engine.Save.State)
			}

			if tc.expectedLog != "" {
				logOutput := strings.TrimSpace(buf.String())
				if !strings.Contains(logOutput, tc.expectedLog) {
					t.Errorf("expected log %q, got %q", tc.expectedLog, logOutput)
				}
			}
		})
	}
}
