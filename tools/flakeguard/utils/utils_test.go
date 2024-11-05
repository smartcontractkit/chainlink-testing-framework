package utils

import (
	"testing"
)

func TestUniqueDirectories(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "unique directories from multiple files",
			input: []string{
				"/path/to/file1.txt",
				"/path/to/file2.txt",
				"/another/path/file3.txt",
			},
			expected: []string{
				"/path/to",
				"/another/path",
			},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name: "same directory multiple files",
			input: []string{
				"/path/to/file1.txt",
				"/path/to/file2.txt",
				"/path/to/file3.txt",
			},
			expected: []string{
				"/path/to",
			},
		},
		{
			name: "mixed directory levels",
			input: []string{
				"/a/b/c/file1.txt",
				"/a/b/file2.txt",
				"/a/file3.txt",
			},
			expected: []string{
				"/a/b/c",
				"/a/b",
				"/a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqueDirectories(tt.input)
			expected := make(map[string]struct{})
			for _, dir := range tt.expected {
				expected[dir] = struct{}{}
			}
			for _, dir := range result {
				if _, found := expected[dir]; !found {
					t.Errorf("unexpected directory: %v", dir)
				}
				delete(expected, dir)
			}
			if len(expected) > 0 {
				t.Errorf("missing directories: %v", expected)
			}
		})
	}
}

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "duplicates",
			input:    []string{"apple", "banana", "apple", "orange", "banana"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			name:     "no duplicates",
			input:    []string{"apple", "banana", "orange"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single item",
			input:    []string{"apple"},
			expected: []string{"apple"},
		},
		{
			name:     "all duplicates",
			input:    []string{"apple", "apple", "apple"},
			expected: []string{"apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Deduplicate(tt.input)
			expected := make(map[string]struct{})
			for _, item := range tt.expected {
				expected[item] = struct{}{}
			}
			for _, item := range result {
				if _, found := expected[item]; !found {
					t.Errorf("unexpected item: %v", item)
				}
				delete(expected, item)
			}
			if len(expected) > 0 {
				t.Errorf("missing items: %v", expected)
			}
		})
	}
}
