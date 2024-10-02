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

var MAX_BYTES_AT_TIME = 100

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

func isBlockSeperator(currentPointer int, bytes []byte) bool {
	if string(bytes[currentPointer:currentPointer+2]) == "\r\n" {
		return true
	}
	return false
}

func readSnapShotFile() {
	f, err := os.Open(SNAPSHOT_FILE_NAME)
	defer f.Close()
	if err != nil {
		zap.L().Error("Error reading snapshot file", zap.Error(err))
		return
	}
	reader := CreateBinaryReader(f)
	n, err := reader.readFromFile()
	if n == 0 {
		zap.L().Error("No contents in file existing")
		return
	}
	if err != nil && err == io.EOF {
		zap.L().Error("End of file found exiting", zap.Error(err))
	}
	err = reader.skipFileHeader()
	if err != nil {
		zap.L().Error("Error skipping file header", zap.Error(err))
	}
	reader.skipBlockSeperator()
	for {
		blockValueType, err := reader.getInt64DataFromBlock()
		if err != nil {
			zap.L().Error("Error while reading snapshot", zap.Error(err))
			return
		}
		fmt.Println("Block value type", blockValueType, "current pointer", reader.currentPointer)
		keyLength, err := reader.getInt64DataFromBlock()
		if err != nil {
			zap.L().Error("Error while reading snapshot", zap.Error(err))
			return
		}
		fmt.Println("Key length", keyLength, "current pointer", reader.currentPointer)
		key, err := reader.getStringDataFromBlock(keyLength)
		if err != nil {
			zap.L().Error("Error while reading snapshot", zap.Error(err))
			return
		}
		fmt.Println("key is", key, reader.currentPointer)
		blockValue, err := reader.getInt64DataFromBlock()
		if err != nil {
			zap.L().Error("Error while reading snapshot", zap.Error(err))
			return
		}
		fmt.Println("block value is", blockValue, reader.currentPointer)
		fmt.Println("----------------------------------------------")
		fmt.Println(key, ":", blockValue)
		fmt.Println("==============================================")
		err = reader.skipBlockSeperator()
		if err != nil {
			zap.L().Error("Error while reading snapshot", zap.Error(err))
			return
		}
	}
}

func showALlbytes() {
	file, _ := os.Open(SNAPSHOT_FILE_NAME)
	bytes := make([]byte, 200)
	file.Read(bytes)
	fmt.Println("All the bytes in the file are", bytes)
	file.Close()
}

func takeSnapShot(wg *sync.WaitGroup, mainMap MainMap) {
	defer wg.Done()
	// file, _ := os.Create(SNAPSHOT_FILE_NAME)
	// value := createBytesForSnapShot(mainMap)
	// value = append(value, byte(END_OF_FILE))
	// fmt.Println(value, "======")
	// file.Write(value)
	// file.Close()
	showALlbytes()
	readSnapShotFile()
}

func RunSnapShotTaker(mainMap MainMap) {
	var wg sync.WaitGroup
	wg.Add(1)
	go takeSnapShot(&wg, mainMap)
	wg.Wait()
}
