package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/letsmakecakes/huffman/pkg/huffman"
)

func main() {
	compress := flag.Bool("c", false, "Compress the input file")
	decompress := flag.Bool("d", false, "Decompress the input file")
	input := flag.String("i", "", "Input file path")
	output := flag.String("o", "", "Output file path")
	flag.Parse()

	if *input == "" {
		fmt.Println("Error: Input file is required")
		flag.Usage()
		os.Exit(1)
	}

	if *output == "" {
		if *compress {
			*output = *input + ".huf"
		} else if *decompress {
			*output = *input + ".dec"
		}
	}

	if *compress && *decompress {
		fmt.Println("Error: Cannot specify both compress and decompress")
		flag.Usage()
		os.Exit(1)
	}

	if *compress {
		if err := huffman.CompressFile(*input, *output); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Compression failed: %v\n", err)
			if err != nil {
				log.Printf("failed to format according to format specifier and write to stderr: %v", err)
			}
			os.Exit(1)
		}

		inputInfo, _ := os.Stat(*input)
		outputInfo, _ := os.Stat(*output)
		ratio := float64(outputInfo.Size()) / float64(inputInfo.Size()) * 100

		fmt.Printf("Compression successful!\n")
		fmt.Printf("Original size: %d bytes\n", inputInfo.Size())
		fmt.Printf("Compressed size: %d bytes\n", outputInfo.Size())
		fmt.Printf("Compression ratio: %.2f%%\n", ratio)
	} else if *decompress {
		if err := huffman.DecompressFile(*input, *output); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Decompression failed: %v\n", err)
			if err != nil {
				log.Printf("failed to format according to format specifier and write to stderr: %v", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Decompression successful! Output written to: %s\n", *output)
	}
}
