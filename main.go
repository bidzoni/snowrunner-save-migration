package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode/utf16"
)

type Record struct {
	Filename string
	Value1   uint32
	Value2   uint16
	Value3   uint16
	Value4   uint64
}

func readContainer(filename string) ([]Record, error) {
	const headerSize = 8
	const recordSize = 160
	const structSize = 16
	const filenameSize = recordSize - 2*structSize

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Skip the header
	_, err = file.Seek(headerSize, 0)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size() - headerSize
	if fileSize%recordSize != 0 {
		return nil, fmt.Errorf("file size is not a multiple of record size")
	}

	recordCount := int(fileSize / recordSize)
	records := make([]Record, 0, recordCount)

	buffer := make([]byte, recordSize)
	for i := 0; i < recordCount; i++ {
		_, err := file.Read(buffer)
		if err != nil {
			return nil, err
		}

		// Read filename (UTF-16LE with null-terminator)
		filenameBytes := buffer[:filenameSize]
		var filenameUTF16 []uint16
		for j := 0; j < len(filenameBytes); j += 2 {
			u := binary.LittleEndian.Uint16(filenameBytes[j : j+2])
			if u == 0 {
				break
			}
			filenameUTF16 = append(filenameUTF16, u)
		}
		// Saves have .cfg extension
		filename := string(utf16.Decode(filenameUTF16)) + ".cfg"

		// Read struct values
		structBytes := buffer[filenameSize : filenameSize+structSize]
		value1 := binary.LittleEndian.Uint32(structBytes[0:4])
		value2 := binary.LittleEndian.Uint16(structBytes[4:6])
		value3 := binary.LittleEndian.Uint16(structBytes[6:8])
		value4 := binary.BigEndian.Uint64(structBytes[8:16])

		records = append(records, Record{
			Filename: filename,
			Value1:   value1,
			Value2:   value2,
			Value3:   value3,
			Value4:   value4,
		})
	}

	return records, nil
}

func createHashCode(record Record) string {
	return fmt.Sprintf("%08X%04X%04X%016X", record.Value1, record.Value2, record.Value3, record.Value4)
}

func migrate(records []Record, sourceDir, destDir string) error {
	for _, record := range records {
		hashCode := createHashCode(record)
		sourcePath := filepath.Join(sourceDir, hashCode)
		destPath := filepath.Join(destDir, record.Filename)

		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", sourcePath, err)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}
		defer destFile.Close()

		fmt.Printf("Copy %s -> %s\n", sourcePath, destPath)
		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("failed to copy from %s to %s: %w", sourcePath, destPath, err)
		}
	}
	return nil
}

func main() {
	containerFile := flag.String("container", "container.51", "Path to the container file")
	sourceDir := flag.String("source", "./", "Path to the source directory containing xbox saves")
	destDir := flag.String("dest", "./cfg", "Path to the destination directory for cfg saves for steam")

	flag.Parse()

	if *containerFile == "" || *sourceDir == "" || *destDir == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	contanerPath := filepath.Join(*sourceDir, *containerFile)
	records, err := readContainer(contanerPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	migrate(records, *sourceDir, *destDir)
}
