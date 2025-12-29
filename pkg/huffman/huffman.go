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

// GenerateCodeTable creates prefix codes from Huffman tree
func GenerateCodeTable(root *Node) CodeTable {
	codes := make(CodeTable)
	if root == nil {
		return codes
	}

	// Special case: single character
	if root.Left == nil && root.Right == nil {
		codes[root.Char] = "0"
		return codes
	}

	generateCodes(root, "", codes)
	return codes
}

func generateCodes(node *Node, code string, codes CodeTable) {
	if node == nil {
		return
	}

	// Leaf node
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = code
		return
	}

	generateCodes(node.Left, code+"0", codes)
	generateCodes(node.Right, code+"1", codes)
}

// EncodeData encodes data using the code table
func EncodeData(data []byte, codes CodeTable) []byte {
	var bitString string
	for _, b := range data {
		bitString += codes[b]
	}

	// Pack bits into bytes
	byteCount := (len(bitString) + 7) / 8
	result := make([]byte, byteCount)

	for i := 0; i < len(bitString); i++ {
		if bitString[i] == '1' {
			byteIdx := i / 8
			bitIdx := 7 - (i % 8)
			result[byteIdx] |= 1 << bitIdx
		}
	}

	return result
}

// DecodeData decodes compressed data using Huffman tree
func DecodeData(data []byte, root *Node, originalSize int64, paddingBits int) ([]byte, error) {
	if root == nil {
		return nil, fmt.Errorf("invalid Huffman tree")
	}

	result := make([]byte, 0, originalSize)
	current := root

	// Special case: single character
	if root.Left == nil && root.Right == nil {
		for i := int64(0); i < originalSize; i++ {
			result = append(result, root.Char)
		}
		return result, nil
	}

	totalBits := len(data)*8 - paddingBits

	for i := 0; i < totalBits && int64(len(result)) < originalSize; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		bit := (data[byteIdx] >> bitIdx) & 1

		if bit == 0 {
			current = current.Left
		} else {
			current = current.Right
		}

		// Reached leaf node
		if current.Left == nil && current.Right == nil {
			result = append(result, current.Char)
			current = root
		}
	}

	return result, nil
}
