package huffman

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// Node represents a node in the Huffman tree.
type Node struct {
	Char  byte
	Freq  int
	Left  *Node
	Right *Node
}

// FrequencyTable stores character frequencies.
type FrequencyTable map[byte]int

// CodeTable stores Huffman codes for each character.
type CodeTable map[byte]string

// BuildFrequencyTable reads a file and counts character occurrences
func BuildFrequencyTable(filename string) (FrequencyTable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(file)

	freq := make(FrequencyTable)
	reader := bufio.NewReader(file)

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		freq[b]++
	}

	if len(freq) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	return freq, nil
}

// BuildFrequencyTableFromData BuildFrequencyTableFrom builds a frequency table from a byte slice
func BuildFrequencyTableFromData(data []byte) FrequencyTable {
	freq := make(FrequencyTable)
	for _, b := range data {
		freq[b]++
	}
	return freq
}
