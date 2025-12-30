package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/letsmakecakes/huffman/pkg/huffman"
)

func TestLargeFileCompression(t *testing.T) {
	// Create a large test file
	testData := bytes.Repeat([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. "), 1000)

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "large_input.txt")
	compressedPath := filepath.Join(tmpDir, "large_compressed.huf")
	decompressedPath := filepath.Join(tmpDir, "large_decompressed.txt")

	// Write test data
	if err := os.WriteFile(inputPath, testData, 0644); err != nil {
		t.Fatal(err)
	}

	// Compress
	if err := huffman.CompressFile(inputPath, compressedPath); err != nil {
		t.Fatalf("Compression failed: %v", err)
	}

	// Verify compressed file exists
	compressedInfo, err := os.Stat(compressedPath)
	if err != nil {
		t.Fatalf("Compressed file not created: %v", err)
	}

	originalInfo, _ := os.Stat(inputPath)
	t.Logf("Original size: %d bytes", originalInfo.Size())
	t.Logf("Compressed size: %d bytes", compressedInfo.Size())
	t.Logf("Compression ratio: %.2f%%", float64(compressedInfo.Size())/float64(originalInfo.Size())*100)

	// Decompress
	if err := huffman.DecompressFile(compressedPath, decompressedPath); err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}

	// Verify content
	decompressed, err := os.ReadFile(decompressedPath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(testData, decompressed) {
		t.Error("Decompressed data doesn't match original")
	}
}

func TestVariousFileTypes(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "text file",
			data: []byte("The quick brown fox jumps over the lazy dog.\n" +
				"Pack my box with five dozen liquor jugs.\n"),
		},
		{
			name: "repetitive data",
			data: bytes.Repeat([]byte("A"), 1000),
		},
		{
			name: "binary-like data",
			data: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			name: "mixed content",
			data: []byte("Hello123!@#\n\t\rWorld456"),
		},
		{
			name: "unicode text",
			data: []byte("Hello, 世界! Привет мир! مرحبا بالعالم"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "input.txt")
			compressedPath := filepath.Join(tmpDir, "compressed.huf")
			decompressedPath := filepath.Join(tmpDir, "decompressed.txt")

			// Write test data
			if err := os.WriteFile(inputPath, tt.data, 0644); err != nil {
				t.Fatal(err)
			}

			// Compress
			if err := huffman.CompressFile(inputPath, compressedPath); err != nil {
				t.Fatalf("Compression failed: %v", err)
			}

			// Decompress
			if err := huffman.DecompressFile(compressedPath, decompressedPath); err != nil {
				t.Fatalf("Decompression failed: %v", err)
			}

			// Verify
			decompressed, err := os.ReadFile(decompressedPath)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(tt.data, decompressed) {
				t.Errorf("Data mismatch for %s", tt.name)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("single byte file", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputPath := filepath.Join(tmpDir, "single.txt")
		compressedPath := filepath.Join(tmpDir, "single.huf")
		decompressedPath := filepath.Join(tmpDir, "single_dec.txt")

		data := []byte("A")
		if err := os.WriteFile(inputPath, data, 0644); err != nil {
			t.Fatal(err)
		}

		if err := huffman.CompressFile(inputPath, compressedPath); err != nil {
			t.Fatalf("Compression failed: %v", err)
		}

		if err := huffman.DecompressFile(compressedPath, decompressedPath); err != nil {
			t.Fatalf("Decompression failed: %v", err)
		}

		decompressed, _ := os.ReadFile(decompressedPath)
		if !bytes.Equal(data, decompressed) {
			t.Error("Single byte file compression/decompression failed")
		}
	})

	t.Run("all same character", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputPath := filepath.Join(tmpDir, "same.txt")
		compressedPath := filepath.Join(tmpDir, "same.huf")
		decompressedPath := filepath.Join(tmpDir, "same_dec.txt")

		data := bytes.Repeat([]byte("X"), 500)
		if err := os.WriteFile(inputPath, data, 0644); err != nil {
			t.Fatal(err)
		}

		if err := huffman.CompressFile(inputPath, compressedPath); err != nil {
			t.Fatalf("Compression failed: %v", err)
		}

		if err := huffman.DecompressFile(compressedPath, decompressedPath); err != nil {
			t.Fatalf("Decompression failed: %v", err)
		}

		decompressed, _ := os.ReadFile(decompressedPath)
		if !bytes.Equal(data, decompressed) {
			t.Error("Same character file compression/decompression failed")
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("non-existent input file", func(t *testing.T) {
		err := huffman.CompressFile("/nonexistent/file.txt", "/tmp/output.huf")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("invalid compressed file", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidPath := filepath.Join(tmpDir, "invalid.huf")
		outputPath := filepath.Join(tmpDir, "output.txt")

		// Write invalid data
		if err := os.WriteFile(invalidPath, []byte("not a valid huffman file"), 0644); err != nil {
			t.Fatal(err)
		}

		err := huffman.DecompressFile(invalidPath, outputPath)
		if err == nil {
			t.Error("Expected error for invalid compressed file")
		}
	})
}

func TestCompressionRatio(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		maxRatio float64 // maximum acceptable compression ratio (lower is better)
	}{
		{
			name:     "highly repetitive",
			data:     bytes.Repeat([]byte("AAAA"), 250),
			maxRatio: 0.20, // should compress to less than 20%
		},
		{
			name:     "english text",
			data:     bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 50),
			maxRatio: 0.70, // should compress to less than 70%
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "input.txt")
			compressedPath := filepath.Join(tmpDir, "compressed.huf")

			if err := os.WriteFile(inputPath, tc.data, 0644); err != nil {
				t.Fatal(err)
			}

			if err := huffman.CompressFile(inputPath, compressedPath); err != nil {
				t.Fatalf("Compression failed: %v", err)
			}

			inputInfo, _ := os.Stat(inputPath)
			compressedInfo, _ := os.Stat(compressedPath)

			ratio := float64(compressedInfo.Size()) / float64(inputInfo.Size())
			t.Logf("Compression ratio: %.2f%%", ratio*100)

			if ratio > tc.maxRatio {
				t.Errorf("Compression ratio %.2f%% exceeds maximum %.2f%%", ratio*100, tc.maxRatio*100)
			}
		})
	}
}
