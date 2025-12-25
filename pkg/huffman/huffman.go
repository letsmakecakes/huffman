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

// BuildHuffmanTree constructs the Huffman tree from frequency table
func BuildHuffmanTree(freq FrequencyTable) *Node {
	if len(freq) == 0 {
		return nil
	}

	// Special case: single character
	if len(freq) == 1 {
		for char, count := range freq {
			return &Node{
				Char:  char,
				Freq:  count,
				Left:  nil,
				Right: nil,
			}
		}
	}

	// Create initial nodes
	nodes := make([]*Node, 0, len(freq))
	for char, count := range freq {
		nodes = append(nodes, &Node{
			Char: char,
			Freq: count,
		})
	}

	// Build tree by repeatedly combining two lowest frequency nodes
	for len(nodes) > 1 {
		// Find two nodes with minimum frequency
		min1Idx, min2Idx := findTwoMinimum(nodes)

		// Create parent node
		parent := &Node{
			Freq:  nodes[min1Idx].Freq + nodes[min2Idx].Freq,
			Left:  nodes[min1Idx],
			Right: nodes[min2Idx],
		}

		// Remove the two minimum nodes and add parent
		nodes = removeNodes(nodes, min1Idx, min2Idx)
		nodes = append(nodes, parent)
	}

	return nodes[0]
}

func findTwoMinimum(nodes []*Node) (int, int) {
	min1, min2 := 0, 1
	if nodes[min1].Freq > nodes[min2].Freq {
		min1, min2 = min2, min1
	}

	for i := 2; i < len(nodes); i++ {
		if nodes[i].Freq < nodes[min1].Freq {
			min2 = min1
			min1 = i
		} else if nodes[i].Freq < nodes[min2].Freq {
			min2 = i
		}
	}

	return min1, min2
}

func removeNodes(nodes []*Node, idx1, idx2 int) []*Node {
	if idx1 > idx2 {
		idx1, idx2 = idx2, idx1
	}
	result := make([]*Node, 0, len(nodes)-2)
	for i, node := range nodes {
		if i != idx1 && i != idx2 {
			result = append(result, node)
		}
	}
	return result
}
