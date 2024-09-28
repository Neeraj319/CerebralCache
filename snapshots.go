package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
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
		if err := binary.Write(&buffer, binary.LittleEndian, int32(STRING_TYPE)); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")

		// write the length of the key
		if err := binary.Write(&buffer, binary.LittleEndian, int32(len(keyBytes))); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")

		// write the key
		if err := binary.Write(&buffer, binary.LittleEndian, keyBytes); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")
		// write the length of the value
		if err := binary.Write(&buffer, binary.LittleEndian, int32(len(valueBytes))); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")
		// write the actual value
		if err := binary.Write(&buffer, binary.LittleEndian, (valueBytes)); err != nil {
			return nil, err
		}
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
		if err := binary.Write(&buffer, binary.LittleEndian, int32(INTEGER_TYPE)); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")

		// write the length of the key
		if err := binary.Write(&buffer, binary.LittleEndian, int32(len(keyBytes))); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")

		// write the key
		if err := binary.Write(&buffer, binary.LittleEndian, keyBytes); err != nil {
			return nil, err
		}
		buffer.WriteString("\r\n")
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
	stringBin, err := convertStringMapToBin(mainMap)
	if err != nil {
		zap.L().Error("Failed creating string map bin", zap.Error(err))
	}
	contentBytes = append(contentBytes, stringBin...)
	integerBin, err := convertIntegerMapToBinary(mainMap)
	if err != nil {
		zap.L().Error("Failed creating integer map bin", zap.Error(err))
	}
	contentBytes = append(contentBytes, integerBin...)
	return contentBytes
}

func readSnapShotFile() {
	f, err := os.Open(SNAPSHOT_FILE_NAME)
	if err != nil {
		zap.L().Error("Error reading snap shot file", zap.Error(err))
	}
	b1 := make([]byte, 1024)
	n1, err := f.Read(b1)
	fmt.Println("content", n1)
}

func takeSnapShot(wg *sync.WaitGroup, mainMap MainMap) {
	defer wg.Done()
	file, _ := os.Create(SNAPSHOT_FILE_NAME)
	value := createBytesForSnapShot(mainMap)
	value = append(value, byte(END_OF_FILE))
	file.Write(value)
	file.Close()
	readSnapShotFile()
}

func RunSnapShotTaker(mainMap MainMap) {
	var wg sync.WaitGroup
	wg.Add(1)
	go takeSnapShot(&wg, mainMap)
	wg.Wait()
}
