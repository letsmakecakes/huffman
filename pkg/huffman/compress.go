package huffman

import (
	"fmt"
	"io"
	"log"
	"os"
)

// CompressFile compresses a file using Huffman encoding
func CompressFile(inputPath, outputPath string) error {
	// Step 1: Build frequency table
	freq, err := BuildFrequencyTable(inputPath)
	if err != nil {
		return fmt.Errorf("failed to build frequency table: %w", err)
	}

	// Step 2: Build a Huffman tree
	tree := BuildHuffmanTree(freq)
	if tree == nil {
		return fmt.Errorf("failed to build huffman tree")
	}

	// Step 3: Generate code table
	codes := GenerateCodeTable(tree)

	// Read original file data
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Step 4: Encode data
	encoded := EncodeData(data, codes)

	// Calculate padding bits
	totalBits := 0
	for _, b := range data {
		totalBits += len(codes[b])
	}
	paddingBits := (8 - (totalBits % 8)) % 8

	// Step 5: Write a compressed file
	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			log.Printf("failed to close output file: %v", err)
		}
	}(output)

	// Write Header
	if err := WriteHeader(output, freq, int64(len(data)), paddingBits); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write encoded data
	if _, err := output.Write(encoded); err != nil {
		return fmt.Errorf("failed to write encoded data: %w", err)
	}

	return nil
}

// DecompressFile decompresses a Huffman encoded file
func DecompressFile(inputPath, outputPath string) error {
	// Open the input file
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func(input *os.File) {
		err := input.Close()
		if err != nil {
			log.Printf("failed to close input file: %v", err)
		}
	}(input)

	// Step 6: Read header
	freq, originalSize, paddingBits, err := ReadHeader(input)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Rebuild Huffman tree
	tree := BuildHuffmanTree(freq)
	if tree == nil {
		return fmt.Errorf("failed to build huffman tree")
	}

	// Step 7: Read and decode compressed data
	encodedData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read encoded data: %w", err)
	}

	decoded, err := DecodeData(encodedData, tree, originalSize, paddingBits)
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	// Write decoded data
	if err := os.WriteFile(outputPath, decoded, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
