package huffman

import (
	"bytes"
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

func TestGenerateCodeTable(t *testing.T) {
	tests := []struct {
		name  string
		freq  FrequencyTable
		check func(CodeTable) bool
	}{
		{
			name: "three characters",
			freq: FrequencyTable{
				'a': 3,
				'b': 2,
				'c': 1,
			},
			check: func(codes CodeTable) bool {
				// More frequent characters should have shorter codes
				return len(codes['a']) <= len(codes['b']) &&
					len(codes['b']) <= len(codes['c']) &&
					len(codes) == 3
			},
		},
		{
			name: "single character",
			freq: FrequencyTable{
				'a': 5,
			},
			check: func(codes CodeTable) bool {
				return codes['a'] == "0" && len(codes) == 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := BuildHuffmanTree(tt.freq)
			codes := GenerateCodeTable(tree)

			if !tt.check(codes) {
				t.Errorf("Code table validation failed: %v", codes)
			}

			// Verify prefix-free property
			if !isPrefixFree(codes) {
				t.Error("Codes are not prefix-free")
			}
		})
	}
}

func isPrefixFree(codes CodeTable) bool {
	codeList := make([]string, 0, len(codes))
	for _, code := range codes {
		codeList = append(codeList, code)
	}

	for i := 0; i < len(codeList); i++ {
		for j := i + 1; j < len(codeList); j++ {
			if isPrefix(codeList[i], codeList[j]) || isPrefix(codeList[j], codeList[i]) {
				return false
			}
		}
	}
	return true
}

func isPrefix(s1, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}
	return s2[:len(s1)] == s1
}

func TestEncodeDecodeData(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "aaabbc"},
		{"single char", "aaaaa"},
		{"longer text", "the quick brown fox jumps over the lazy dog"},
		{"with newlines", "hello\nworld\n"},
		{"unicode", "Hello, 世界!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(tt.input)
			freq := BuildFrequencyTableFromData(data)
			tree := BuildHuffmanTree(freq)
			codes := GenerateCodeTable(tree)

			// Encode
			encoded := EncodeData(data, codes)

			// Calculate padding
			totalBits := 0
			for _, b := range data {
				totalBits += len(codes[b])
			}
			paddingBits := (8 - (totalBits % 8)) % 8

			// Decode
			decoded, err := DecodeData(encoded, tree, int64(len(data)), paddingBits)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}

			if !bytes.Equal(data, decoded) {
				t.Errorf("Decoded data doesn't match original.\nOriginal: %s\nDecoded: %s", data, decoded)
			}
		})
	}
}

func TestHeaderWriteRead(t *testing.T) {
	freq := FrequencyTable{
		'a': 3,
		'b': 2,
		'c': 1,
	}
	originalSize := int64(100)
	paddingBits := 5

	var buf bytes.Buffer

	// Write header
	err := WriteHeader(&buf, freq, originalSize, paddingBits)
	if err != nil {
		t.Fatalf("WriteHeader error: %v", err)
	}

	// Read header
	readFreq, readSize, readPadding, err := ReadHeader(&buf)
	if err != nil {
		t.Fatalf("ReadHeader error: %v", err)
	}

	if !reflect.DeepEqual(freq, readFreq) {
		t.Errorf("Frequency tables don't match.\nExpected: %v\nGot: %v", freq, readFreq)
	}

	if originalSize != readSize {
		t.Errorf("Original sizes don't match. Expected: %d, Got: %d", originalSize, readSize)
	}

	if paddingBits != readPadding {
		t.Errorf("Padding bits don't match. Expected: %d, Got: %d", paddingBits, readPadding)
	}
}
