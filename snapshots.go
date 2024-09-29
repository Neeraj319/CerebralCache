package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
)

func createFileHeader() ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString(FILE_HEADER); err != nil {
		return nil, err
	}

	if _, err := buffer.WriteString("\r\n"); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func convertStringMapToBin(mainMap MainMap) ([]byte, error) {
	stringMap := mainMap.STRING_MAP
	var buffer bytes.Buffer
	for key, value := range stringMap {
		keyBytes := []byte(key)
		valueBytes := []byte(value)

		// write the type of the value
		if err := binary.Write(&buffer, binary.LittleEndian, int64(STRING_TYPE)); err != nil {
			return nil, err
		}
		// write the length of the key
		if err := binary.Write(&buffer, binary.LittleEndian, int64(len(keyBytes))); err != nil {
			return nil, err
		}
		// write the key
		if err := binary.Write(&buffer, binary.LittleEndian, keyBytes); err != nil {
			return nil, err
		}
		buffer.WriteByte(byte(0))
		// write the length of the value
		if err := binary.Write(&buffer, binary.LittleEndian, int64(len(valueBytes))); err != nil {
			return nil, err
		}
		// write the actual value
		if err := binary.Write(&buffer, binary.LittleEndian, (valueBytes)); err != nil {
			return nil, err
		}
		buffer.WriteByte(byte(0))
		buffer.WriteString("\r\n")
	}
	return buffer.Bytes(), nil
}

func convertIntegerMapToBinary(minmap MainMap) ([]byte, error) {
	var buffer bytes.Buffer
	integerMap := minmap.INTEGER_MAP
	for key, value := range integerMap {
		keyBytes := []byte(key)

		// write the type of the value
		if err := binary.Write(&buffer, binary.LittleEndian, int64(INTEGER_TYPE)); err != nil {
			return nil, err
		}
		// write the length of the key
		if err := binary.Write(&buffer, binary.LittleEndian, int64(len(keyBytes))); err != nil {
			return nil, err
		}
		// write the key
		if err := binary.Write(&buffer, binary.LittleEndian, keyBytes); err != nil {
			return nil, err
		}
		buffer.WriteByte(byte(0))
		// write the actual value
		if err := binary.Write(&buffer, binary.LittleEndian, (value)); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")
	}
	return buffer.Bytes(), nil
}

func createBytesForSnapShot(mainMap MainMap) []byte {
	contentBytes, err := createFileHeader()
	if err != nil {
		zap.L().Error("Failed to create file header", zap.Error(err))
	}
	integerBin, err := convertIntegerMapToBinary(mainMap)
	if err != nil {
		zap.L().Error("Failed creating integer map bin", zap.Error(err))
	}
	contentBytes = append(contentBytes, integerBin...)
	stringBin, err := convertStringMapToBin(mainMap)
	if err != nil {
		zap.L().Error("Failed creating string map bin", zap.Error(err))
	}
	contentBytes = append(contentBytes, stringBin...)
	return contentBytes
}

func processReadBytes(bytes []byte) {
	fmt.Println("bytes to process are", bytes)
	indx := 0
	dataType := binary.LittleEndian.Uint64(bytes[indx : indx+8])
	fmt.Println("type of the data is", dataType)
	indx += 8
	keyLength := binary.LittleEndian.Uint64(bytes[indx : indx+8])
	fmt.Println("Length of the key is", keyLength)
	keyBytes := make([]byte, keyLength)
	indx += 8
	current := 0
	fmt.Println("index is", indx)
	for {
		if bytes[indx] == 0 {
			break
		}
		keyBytes[current] = bytes[indx]
		current++
		indx++
	}
	fmt.Println("The key is", string(keyBytes))
	indx++
	value := binary.LittleEndian.Uint64(bytes[indx:])
	fmt.Println("The valu eis", value)
}

func skipFileHeader(bytes []byte) int {
	if string(bytes[0:6]) == FILE_HEADER {
		return 6
	} else {
		zap.L().Error("Invalid file header exiting")
		return 0
	}
}
func isBlockSeperator(currentPointer int, bytes []byte) bool {
	if string(bytes[currentPointer:currentPointer+2]) == "\r\n" {
		return true
	}
	return false
}
func skipBlockSeperator(currentPointer int, bytes []byte) int {
	if isBlockSeperator(currentPointer, bytes) {
		return currentPointer + 2
	}
	return currentPointer
}

func findNextBlock(currentPointer int, bytes []byte) int {
	for {
		if (currentPointer + 2) >= len(bytes) {
			return 0
		}
		if isBlockSeperator(currentPointer, bytes) {
			return currentPointer
		}
		currentPointer += 2
	}
}

func readSnapShotFile() {
	f, err := os.Open(SNAPSHOT_FILE_NAME)
	defer f.Close()
	if err != nil {
		zap.L().Error("Error reading snap shot file", zap.Error(err))
		return
	}
	var buffer bytes.Buffer
	bytes := make([]byte, 100)
	currentPointer := 0
	l, err := f.Read(bytes)
	if err == io.EOF {
		return
	}
	// fmt.Println("Current bytes before header skip", bytes[currentPointer:l])
	currentPointer = skipFileHeader(bytes)
	// fmt.Println("Current bytes after header skip", bytes[currentPointer:l])
	for {
		currentPointer = skipBlockSeperator(currentPointer, bytes)
		fmt.Println("Current bytes after block skip", bytes[currentPointer:l])
		nextBlockIndex := findNextBlock(currentPointer, bytes)
		fmt.Println("the next pointer is", nextBlockIndex)
		if nextBlockIndex == 0 {
			fmt.Println("You are here and i know it")
			buffer.Write(bytes[currentPointer:])
			fmt.Println("THe new buffer is", buffer.Bytes())
			currentPointer = 0
			nextBlockIndex = len(bytes)
		} else {
			fmt.Println("Here nex block is at", nextBlockIndex, "current pointer is at", currentPointer)
			buffer.Write(bytes[currentPointer:nextBlockIndex])
			processReadBytes(buffer.Bytes())
			fmt.Println("=========================")
			currentPointer = nextBlockIndex
			fmt.Println("Current pointer is here", currentPointer)
			buffer.Reset()
		}
		if nextBlockIndex < len(bytes) {
			continue
		}
		_, err = f.Read(bytes)
		if err == io.EOF {
			fmt.Println("Here hehehe")
			return
		}
	}
}

func takeSnapShot(wg *sync.WaitGroup, mainMap MainMap) {
	defer wg.Done()
	// file, _ := os.Create(SNAPSHOT_FILE_NAME)
	// file.Close()
	// value := createBytesForSnapShot(mainMap)
	// value = append(value, byte(END_OF_FILE))
	// fmt.Println("the length of the file is", len(value))
	// file.Write(value)
	// file.Close()
	readSnapShotFile()
}

func RunSnapShotTaker(mainMap MainMap) {
	var wg sync.WaitGroup
	wg.Add(1)
	go takeSnapShot(&wg, mainMap)
	wg.Wait()
}
