package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type BinaryReader struct {
	file               *os.File
	bytesAtATime       []byte
	currentPointer     int
	currentBytesLength int
}

func CreateBinaryReader(file *os.File) BinaryReader {
	return BinaryReader{
		file:         file,
		bytesAtATime: make([]byte, MAX_BYTES_AT_TIME),
	}
}
func (reader *BinaryReader) readFromFile() (int, error) {
	n, err := reader.file.Read(reader.bytesAtATime)
	reader.currentBytesLength = n
	return n, err
}

func (reader *BinaryReader) readWithOtherByteArray(bytes []byte) (int, error) {
	l, err := reader.file.Read(bytes)
	reader.currentBytesLength = l
	if err != nil {
		return 0, err
	}
	return l, nil
}

func (reader *BinaryReader) skipFileHeader() error {
	if string(reader.bytesAtATime[0:len(FILE_HEADER)]) == FILE_HEADER {
		reader.currentPointer += len(FILE_HEADER)
		return nil
	}
	return fmt.Errorf("Invalid file header")
}

func (reader *BinaryReader) skipBlockSeperator() error {
	if reader.willPointerGoBeyondLimit(2) {
		err := reader.readMoreBytes()
		return err
	}
	if isBlockSeperator(reader.currentPointer, reader.bytesAtATime) {
		reader.currentPointer += 2
		return nil
	}
	return nil
}
func (reader *BinaryReader) willPointerGoBeyondLimit(valueToIncrease int) bool {
	if (reader.currentPointer + valueToIncrease) > reader.currentBytesLength {
		return true
	}
	return false
}
func (reader *BinaryReader) resetCurrentPointer() {
	reader.currentPointer = 0
}

func (reader *BinaryReader) readMoreBytes() error {
	_, error := reader.readFromFile()
	if error != nil {
		return error
	}
	toSkipAt := reader.currentPointer - reader.currentBytesLength
	reader.resetCurrentPointer()
	reader.increaseCurrentPointer(toSkipAt)
	return nil
}
func (reader *BinaryReader) increaseCurrentPointer(increaseTo int) {
	reader.currentPointer += increaseTo
}

func (reader *BinaryReader) getInt64DataFromBlock() (int64, error) {
	if reader.willPointerGoBeyondLimit(8) {
		fmt.Println("Pointer goes beyond the limit ofc")
		remaningBytes := make([]byte, MAX_BYTES_AT_TIME)
		transferBytes(reader.bytesAtATime[reader.currentPointer:], remaningBytes, 0)
		bytesToReadFromFile := reader.currentPointer
		newBytes := make([]byte, bytesToReadFromFile)
		l, err := reader.readWithOtherByteArray(newBytes)
		if l < 8 {
			return 0, fmt.Errorf("String length didn't match with the specified length, expected: %d got: %d", 8, l)
		}
		if err != nil {
			return 0, err
		}
		transferBytes(reader.bytesAtATime, remaningBytes, (reader.currentBytesLength-reader.currentPointer)-1)
		// fmt.Println("new bytes that were read are", newBytes, "previous bytes are", remaningBytes)
		// err := reader.readMoreBytes()
		// fmt.Println("Read more bytes", reader.bytesAtATime, reader.currentPointer)
		// if err != nil {
		// 	return 0, err
		// }
	}
	from := reader.currentPointer
	reader.increaseCurrentPointer(8)
	to := reader.currentPointer
	return int64(binary.LittleEndian.Uint64(reader.bytesAtATime[from:to])), nil
}
func transferBytes(fromBytes []byte, toBytes []byte, indx int) {
	for _, value := range fromBytes {
		toBytes[indx] = value
		indx++
	}
}

func (reader *BinaryReader) getStringDataFromBlock(stringLength int64) (string, error) {
	bytes := make([]byte, stringLength)
	if reader.willPointerGoBeyondLimit(int(stringLength)) {
		// fmt.Println("Pointer goes beyond the limit wawa")
		transferBytes(reader.bytesAtATime[reader.currentPointer:], bytes, 0)
		// fmt.Println("While string pointer goes beyond limit now the bytes are", bytes)
		spaceRemaning := (reader.currentPointer + int(stringLength)) - reader.currentBytesLength
		remaningByteToRead := make([]byte, spaceRemaning)
		l, err := reader.readWithOtherByteArray(remaningByteToRead)
		// fmt.Println("Here after reading", remaningByteToRead, bytes, spaceRemaning)
		if l < spaceRemaning {
			return "", fmt.Errorf("String length didn't match with the specified length, expected: %d got: %d", stringLength, l)
		}
		transferBytes(remaningByteToRead, bytes, spaceRemaning-1)
		if err != nil {
			return "", nil
		}
		reader.increaseCurrentPointer(spaceRemaning)
		return string(bytes), nil
	} else {
		from := reader.currentPointer
		bytes = append(bytes, reader.bytesAtATime[from:from+int(stringLength)]...)
		reader.increaseCurrentPointer(int(stringLength) + 1)
		return string(bytes), nil
	}
}
