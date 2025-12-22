/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package jsonpath

import (
	"encoding/json"
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		path     string
		expected []PathMatch
		wantErr  bool
	}{
		{
			name: "simple field",
			json: `{"email": "test@example.com"}`,
			path: "$.email",
			expected: []PathMatch{
				{Path: "$.email", Value: "test@example.com"},
			},
		},
		{
			name: "nested field",
			json: `{"user": {"email": "test@example.com"}}`,
			path: "$.user.email",
			expected: []PathMatch{
				{Path: "$.user.email", Value: "test@example.com"},
			},
		},
		{
			name: "array wildcard",
			json: `{"users": [{"email": "a@test.com"}, {"email": "b@test.com"}]}`,
			path: "$.users[*].email",
			expected: []PathMatch{
				{Path: "$.users[0].email", Value: "a@test.com"},
				{Path: "$.users[1].email", Value: "b@test.com"},
			},
		},
		{
			name: "specific array index",
			json: `{"users": [{"email": "a@test.com"}, {"email": "b@test.com"}]}`,
			path: "$.users[1].email",
			expected: []PathMatch{
				{Path: "$.users[1].email", Value: "b@test.com"},
			},
		},
		{
			name:     "missing path",
			json:     `{"name": "John"}`,
			path:     "$.email",
			expected: nil,
		},
		{
			name:     "null value",
			json:     `{"email": null}`,
			path:     "$.email",
			expected: nil, // null values are skipped
		},
		{
			name:    "invalid json",
			json:    `not json`,
			path:    "$.email",
			wantErr: true,
		},
		{
			name:    "invalid path",
			json:    `{"email": "test"}`,
			path:    "invalid path",
			wantErr: true,
		},
		{
			name: "simple array values",
			json: `{"tags": ["tag1", "tag2", "tag3"]}`,
			path: "$.tags[*]",
			expected: []PathMatch{
				{Path: "$.tags[0]", Value: "tag1"},
				{Path: "$.tags[1]", Value: "tag2"},
				{Path: "$.tags[2]", Value: "tag3"},
			},
		},
	}

	p := NewProcessor(true) // quiet mode for tests

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := p.Extract([]byte(tt.json), tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(matches) != len(tt.expected) {
				t.Errorf("got %d matches, want %d", len(matches), len(tt.expected))
				return
			}

			for i, match := range matches {
				if match.Path != tt.expected[i].Path {
					t.Errorf("match[%d].Path = %q, want %q", i, match.Path, tt.expected[i].Path)
				}
				if match.Value != tt.expected[i].Value {
					t.Errorf("match[%d].Value = %q, want %q", i, match.Value, tt.expected[i].Value)
				}
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name         string
		json         string
		replacements map[string]string
		wantJSON     string
		wantErr      bool
	}{
		{
			name: "simple field",
			json: `{"email": "old@test.com"}`,
			replacements: map[string]string{
				"$.email": "new@test.com",
			},
			wantJSON: `{"email":"new@test.com"}`,
		},
		{
			name: "nested field",
			json: `{"user": {"email": "old@test.com", "name": "John"}}`,
			replacements: map[string]string{
				"$.user.email": "new@test.com",
			},
			wantJSON: `{"user":{"email":"new@test.com","name":"John"}}`,
		},
		{
			name: "multiple replacements",
			json: `{"email": "a@test.com", "phone": "123"}`,
			replacements: map[string]string{
				"$.email": "b@test.com",
				"$.phone": "456",
			},
			wantJSON: `{"email":"b@test.com","phone":"456"}`,
		},
		{
			name: "array element",
			json: `{"users": [{"email": "a@test.com"}, {"email": "b@test.com"}]}`,
			replacements: map[string]string{
				"$.users[0].email": "x@test.com",
				"$.users[1].email": "y@test.com",
			},
			wantJSON: `{"users":[{"email":"x@test.com"},{"email":"y@test.com"}]}`,
		},
		{
			name:         "no replacements",
			json:         `{"email": "test@test.com"}`,
			replacements: map[string]string{},
			wantJSON:     `{"email": "test@test.com"}`,
		},
		{
			name:    "invalid json",
			json:    `not json`,
			wantErr: true,
			replacements: map[string]string{
				"$.email": "test",
			},
		},
	}

	p := NewProcessor(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := p.Replace([]byte(tt.json), tt.replacements)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Normalize JSON for comparison
			var gotObj, wantObj interface{}
			if err := json.Unmarshal(result, &gotObj); err != nil {
				t.Errorf("result is not valid JSON: %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.wantJSON), &wantObj); err != nil {
				t.Errorf("wantJSON is not valid JSON: %v", err)
				return
			}

			gotNorm, _ := json.Marshal(gotObj)
			wantNorm, _ := json.Marshal(wantObj)

			if string(gotNorm) != string(wantNorm) {
				t.Errorf("got %s, want %s", string(gotNorm), string(wantNorm))
			}
		})
	}
}

func TestExtractAndCollect(t *testing.T) {
	json := `{
        "email": "test@example.com",
        "phone": "123-456-7890",
        "contacts": [
            {"name": "John", "email": "john@test.com"},
            {"name": "Jane", "email": "jane@test.com"}
        ]
    }`

	p := NewProcessor(true)

	paths := []string{"$.email", "$.phone", "$.contacts[*].email"}

	result, err := p.ExtractAndCollect([]byte(json), paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check $.email
	if len(result["$.email"]) != 1 {
		t.Errorf("expected 1 match for $.email, got %d", len(result["$.email"]))
	} else if result["$.email"][0].Value != "test@example.com" {
		t.Errorf("$.email value = %q, want %q", result["$.email"][0].Value, "test@example.com")
	}

	// Check $.phone
	if len(result["$.phone"]) != 1 {
		t.Errorf("expected 1 match for $.phone, got %d", len(result["$.phone"]))
	}

	// Check $.contacts[*].email
	if len(result["$.contacts[*].email"]) != 2 {
		t.Errorf("expected 2 matches for $.contacts[*].email, got %d",
			len(result["$.contacts[*].email"]))
	}
}

func TestBuildConcretePath(t *testing.T) {
	tests := []struct {
		pathExpr string
		index    int
		total    int
		want     string
	}{
		{"$.email", 0, 1, "$.email"},
		{"$.users[*].email", 0, 2, "$.users[0].email"},
		{"$.users[*].email", 1, 2, "$.users[1].email"},
		{"$.data[*]", 5, 10, "$.data[5]"},
		{"$.a[*].b[*].c", 0, 3, "$.a[0].b[*].c"}, // Only first wildcard replaced
	}

	for _, tt := range tests {
		t.Run(tt.pathExpr, func(t *testing.T) {
			got := buildConcretePath(tt.pathExpr, tt.index, tt.total)
			if got != tt.want {
				t.Errorf("buildConcretePath(%q, %d, %d) = %q, want %q",
					tt.pathExpr, tt.index, tt.total, got, tt.want)
			}
		})
	}
}
