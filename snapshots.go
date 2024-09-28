package main

import (
	"bytes"
	"encoding/binary"
	"go.uber.org/zap"
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

func createBytesForSnapShot(mainMap MainMap) []byte {
	fileHeader, err := createFileHeader()
	if err != nil {
		zap.L().Error("Failed to create file header", zap.Error(err))
	}
	stringBin, err := convertStringMapToBin(mainMap)
	if err != nil {
		zap.L().Error("Failed creating string map bin", zap.Error(err))
	}
	fileHeader = append(fileHeader, stringBin...)
	return fileHeader
}

func takeSnapShot(wg *sync.WaitGroup, mainMap MainMap) {
	defer wg.Done()
	file, _ := os.Create("something.bin")
	value := createBytesForSnapShot(mainMap)
	file.Write(value)
	file.Close()
}

func RunSnapShotTaker(mainMap MainMap) {
	var wg sync.WaitGroup
	wg.Add(1)
	go takeSnapShot(&wg, mainMap)
	wg.Wait()
}
