package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type BinaryReader struct {
	file *os.File
}

func CreateBinaryReader(file *os.File) BinaryReader {
	return BinaryReader{
		file: file,
	}
}

func (reader *BinaryReader) skipFileHeader() error {
	bytes := make([]byte, len(FILE_HEADER))
	l, err := reader.file.Read(bytes)
	if err != nil {
		return err
	}
	if l < len(FILE_HEADER) {
		return fmt.Errorf("Error while skipping file header expected %d bytes found %d", len(FILE_HEADER), l)
	}
	header_contents := string(bytes)
	if header_contents != FILE_HEADER {
		return fmt.Errorf("Invalid file header found, expected %s, found %s", FILE_HEADER, header_contents)
	}
	return err
}

func (reader *BinaryReader) skipBlockSeperator() error {
	bytes := make([]byte, BLOCK_SEPERATOR_LENGTH)
	l, err := reader.file.Read(bytes)
	if err != nil {
		return err
	}
	if l < BLOCK_SEPERATOR_LENGTH {
		return fmt.Errorf("Error skipping block seperator not enough bytes found")
	}
	blockSeperatorString := string(bytes)
	if blockSeperatorString != "\r\n" {
		return fmt.Errorf("Block seperator expected '\r\n' found %s", blockSeperatorString)
	}
	return nil
}

func (reader *BinaryReader) getInt64DataFromBlock() (int64, error) {
	bytes := make([]byte, INT_TYPE_LENGTH)
	l, err := reader.file.Read(bytes)
	if err != nil {
		return 0, err
	}
	if l < INT_TYPE_LENGTH {
		return 0, fmt.Errorf("Expected 8 bytes but found only %d while reading integer", l)
	}

	return int64(binary.LittleEndian.Uint64(bytes)), nil
}

func (reader *BinaryReader) getStringDataFromBlock(stringLength int64) (string, error) {
	bytes := make([]byte, stringLength+1)
	l, err := reader.file.Read(bytes)
	if err != nil {
		return "", err
	}
	if l < int(stringLength+1) {
		return "", fmt.Errorf("Expected %d bytes but found only %d while reading string", stringLength, l)
	}
	return string(bytes[:stringLength]), err
}
