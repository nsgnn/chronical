package main

import (
	"testing"
)

func TestLoadLevel(t *testing.T) {
	testCases := []struct {
		name             string
		id               int
		levelName        string
		author           string
		solution         string
		expectError      bool
		expectedSolution string
	}{
		{
			name:             "Valid 4x4 Level",
			id:               1,
			levelName:        "Standard",
			author:           "Admin",
			solution:         "1.2.\n.3.4\n2.4.\n.1.3",
			expectError:      false,
			expectedSolution: "1.2.\n.3.4\n2.4.\n.1.3",
		},
		{
			name:             "Valid Rectangular Level",
			id:               2,
			levelName:        "Rect",
			author:           "Admin",
			solution:         "123\n456",
			expectError:      false,
			expectedSolution: "123\n456",
		},
		{
			name:             "Valid Level with Trailing Spaces",
			id:               3,
			levelName:        "Trailing Space",
			author:           "Admin",
			solution:         "1. \n.2 \n",
			expectError:      false,
			expectedSolution: "1. \n.2 ",
		},
		{
			name:        "Invalid Negative ID",
			id:          -1,
			levelName:   "Bad ID",
			author:      "Admin",
			solution:    "12\n34",
			expectError: true,
		},
		{
			name:        "Empty Solution",
			id:          4,
			levelName:   "Empty",
			author:      "Admin",
			solution:    "",
			expectError: true,
		},
		{
			name:        "Solution with Uneven Rows (Oddly Shaped)",
			id:          5,
			levelName:   "Odd Shape",
			author:      "Admin",
			solution:    "123\n45\n678",
			expectError: true,
		},
		{
			name:        "Solution with Empty Row (Malformed)",
			id:          6,
			levelName:   "Malformed",
			author:      "Admin",
			solution:    "12\n\n34",
			expectError: true,
		},
		{
			name:        "Solution with only a newline",
			id:          7,
			levelName:   "Newline only",
			author:      "Admin",
			solution:    "\n",
			expectError: true,
		},
		{
			name:             "Solution with spaces and tabs",
			id:               8,
			levelName:        "Whitespace",
			author:           "Admin",
			solution:         "1 2\t\n3 4\t",
			expectError:      false,
			expectedSolution: "1 2\t\n3 4\t",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			level, err := LoadLevel(tc.id, tc.levelName, tc.author, tc.solution)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if level == nil {
					t.Errorf("Expected a level object, but got nil")
					return // Avoid panic on nil dereference
				}
				if level.id != tc.id {
					t.Errorf("Expected level ID %d, but got %d", tc.id, level.id)
				}
				if level.solution != tc.expectedSolution {
					t.Errorf("Expected solution '%s', but got '%s'", tc.expectedSolution, level.solution)
				}
			}
		})
	}
}
