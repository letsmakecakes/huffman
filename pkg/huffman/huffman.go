package huffman

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Node represents a node in the Huffman tree.
type Node struct {
	Char  byte
	Freq  int
	Seq   int // Sequence number for tie-breaking in tree building
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

// BuildHuffmanTree constructs the Huffman tree from a frequency table
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
				Seq:   0,
				Left:  nil,
				Right: nil,
			}
		}
	}

	// Create initial nodes sorted by character for deterministic ordering
	// This ensures that regardless of map iteration order, we always get the same tree
	nodes := make([]*Node, 0, len(freq))
	chars := make([]byte, 0, len(freq))
	for char := range freq {
		chars = append(chars, char)
	}
	// Simple bubble sort for small character sets (0-255)
	for i := 0; i < len(chars); i++ {
		for j := i + 1; j < len(chars); j++ {
			if chars[j] < chars[i] {
				chars[i], chars[j] = chars[j], chars[i]
			}
		}
	}

	seq := 0
	for _, char := range chars {
		nodes = append(nodes, &Node{
			Char: char,
			Freq: freq[char],
			Seq:  seq,
		})
		seq++
	}

	// Build tree by repeatedly combining two lowest frequency nodes
	for len(nodes) > 1 {
		// Find two nodes with a minimum frequency
		min1Idx, min2Idx := findTwoMinimum(nodes)

		// Create parent node
		parent := &Node{
			Freq:  nodes[min1Idx].Freq + nodes[min2Idx].Freq,
			Seq:   seq,
			Left:  nodes[min1Idx],
			Right: nodes[min2Idx],
		}
		seq++

		// Remove the two minimum nodes and add a parent
		nodes = removeNodes(nodes, min1Idx, min2Idx)
		nodes = append(nodes, parent)
	}

	return nodes[0]
}

func findTwoMinimum(nodes []*Node) (int, int) {
	min1, min2 := 0, 1
	if nodes[min1].Freq > nodes[min2].Freq ||
		(nodes[min1].Freq == nodes[min2].Freq && nodes[min1].Seq > nodes[min2].Seq) {
		min1, min2 = min2, min1
	}

	for i := 2; i < len(nodes); i++ {
		if nodes[i].Freq < nodes[min1].Freq ||
			(nodes[i].Freq == nodes[min1].Freq && nodes[i].Seq < nodes[min1].Seq) {
			min2 = min1
			min1 = i
		} else if nodes[i].Freq < nodes[min2].Freq ||
			(nodes[i].Freq == nodes[min2].Freq && nodes[i].Seq < nodes[min2].Seq) {
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

// GenerateCodeTable creates prefix codes from a Huffman tree
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
	// Use strings.Builder for efficient string concatenation
	var buf strings.Builder
	for _, b := range data {
		buf.WriteString(codes[b])
	}
	bitString := buf.String()

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

// WriteHeader writes a compression header to an output file
func WriteHeader(writer io.Writer, freq FrequencyTable, originalSize int64, paddingBits int) error {
	// Write a magic byte
	if err := binary.Write(writer, binary.BigEndian, uint8(0x48)); err != nil { // 'H'
		return err
	}

	// Write original file size as uint32
	if err := binary.Write(writer, binary.BigEndian, uint32(originalSize)); err != nil {
		return err
	}

	// Write padding bits and table size (1 byte)
	tableSize := uint8(len(freq))
	paddingByte := (uint8(paddingBits) << 5) | (tableSize & 0x1F)
	if err := binary.Write(writer, binary.BigEndian, paddingByte); err != nil {
		return err
	}

	// Write character and frequency for each entry
	for char, count := range freq {
		if err := binary.Write(writer, binary.BigEndian, char); err != nil {
			return err
		}
		if err := binary.Write(writer, binary.BigEndian, uint8(count)); err != nil {
			return err
		}
	}

	return nil
}

// ReadHeader reads compression header from an input file
func ReadHeader(reader io.Reader) (FrequencyTable, int64, int, error) {
	// Read and verify the magic byte
	var magic uint8
	if err := binary.Read(reader, binary.BigEndian, &magic); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read magic byte: %w", err)
	}
	if magic != 0x48 { // 'H'
		return nil, 0, 0, fmt.Errorf("invalid file format")
	}

	// Read the original file size as uint32
	var originalSize uint32
	if err := binary.Read(reader, binary.BigEndian, &originalSize); err != nil {
		return nil, 0, 0, err
	}

	// Read padding bits and table size from one byte
	var paddingAndSize uint8
	if err := binary.Read(reader, binary.BigEndian, &paddingAndSize); err != nil {
		return nil, 0, 0, err
	}
	paddingBits := int(paddingAndSize >> 5)
	tableSize := paddingAndSize & 0x1F

	// Read the frequency table
	freq := make(FrequencyTable)
	for i := uint8(0); i < tableSize; i++ {
		var char byte
		if err := binary.Read(reader, binary.BigEndian, &char); err != nil {
			return nil, 0, 0, err
		}

		var count uint8
		if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
			return nil, 0, 0, err
		}

		freq[char] = int(count)
	}

	return freq, int64(originalSize), paddingBits, nil
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
			if current.Left == nil {
				return nil, fmt.Errorf("invalid bit sequence: no left child at bit %d", i)
			}
			current = current.Left
		} else {
			if current.Right == nil {
				return nil, fmt.Errorf("invalid bit sequence: no right child at bit %d", i)
			}
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
