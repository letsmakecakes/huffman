package huffman

import (
	"reflect"
	"testing"
)

func TestBuildFrequencyTable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected FrequencyTable
	}{
		{
			name:  "simple string",
			input: "aaabbc",
			expected: FrequencyTable{
				'a': 3,
				'b': 2,
				'c': 1,
			},
		},
		{
			name:  "single character",
			input: "aaaaa",
			expected: FrequencyTable{
				'a': 5,
			},
		},
		{
			name:  "all unique",
			input: "abcdef",
			expected: FrequencyTable{
				'a': 1,
				'b': 1,
				'c': 1,
				'd': 1,
				'e': 1,
				'f': 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildFrequencyTableFromData([]byte(tt.input))
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBuildHuffmanTree(t *testing.T) {
	freq := FrequencyTable{
		'a': 3,
		'b': 2,
		'c': 1,
	}

	tree := BuildHuffmanTree(freq)

	if tree == nil {
		t.Fatal("Expected non-nil tree")
	}

	if tree.Freq != 6 {
		t.Errorf("Expected root frequency 6, got %d", tree.Freq)
	}

	// Tree should have children
	if tree.Left == nil || tree.Right == nil {
		t.Errorf("Expected tree to have both children")
	}
}

func TestBuildHuffmanTreeSingleChar(t *testing.T) {
	freq := FrequencyTable{
		'a': 5,
	}

	tree := BuildHuffmanTree(freq)

	if tree == nil {
		t.Fatal("Expected non-nil tree")
	}

	if tree.Char != 'a' || tree.Freq != 5 {
		t.Errorf("Expected char 'a with freq 5, got char '%c' with freq %d", tree.Char, tree.Freq)
	}
}
