package snapshots

import (
	"bytes"
	"encoding/binary"
	"in-memory-store/constants"
	"in-memory-store/schemas"
	"os"
	"sync"

	"go.uber.org/zap"
)

func createFileHeader() (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString(constants.FILE_HEADER); err != nil {
		return nil, err
	}

	if err := binary.Write(&buffer, binary.LittleEndian, int64(constants.CURRENT_VERSION)); err != nil {
		return nil, err
	}
	return &buffer, nil
}

func convertStringMapToBin(mainMap schemas.MainMap) ([]byte, error) {
	stringMap := mainMap.STRING_MAP
	var buffer bytes.Buffer
	for key, value := range stringMap {
		keyBytes := []byte(key)
		valueBytes := []byte(value)

		// write the type of the value
		if err := binary.Write(&buffer, binary.LittleEndian, int64(constants.STRING_TYPE)); err != nil {
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

func convertIntegerMapToBinary(minmap schemas.MainMap) ([]byte, error) {
	var buffer bytes.Buffer
	integerMap := minmap.INTEGER_MAP
	for key, value := range integerMap {
		keyBytes := []byte(key)

		// write the type of the value
		if err := binary.Write(&buffer, binary.LittleEndian, int64(constants.INTEGER_TYPE)); err != nil {
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

func createBytesForSnapShot(mainMap schemas.MainMap) *bytes.Buffer {
	mainBuffer, err := createFileHeader()
	if err != nil {
		zap.L().Error("Failed to create file header", zap.Error(err))
	}
	integerBinBytes, err := convertIntegerMapToBinary(mainMap)
	if err != nil {
		zap.L().Error("Failed creating integer map bin", zap.Error(err))
	}
	mainBuffer.Write(integerBinBytes)
	stringBinBytes, err := convertStringMapToBin(mainMap)
	if err != nil {
		zap.L().Error("Failed creating string map bin", zap.Error(err))
	}
	mainBuffer.Write(stringBinBytes)
	return mainBuffer
}
func takeSnapShot(wg *sync.WaitGroup, mainMap schemas.MainMap) {
	defer wg.Done()
	file, err := os.Create(constants.SNAPSHOT_FILE_NAME)
	if err != nil {
		zap.L().Error("Error while taking snapshot of file", zap.Error(err))
	}
	buffer := createBytesForSnapShot(mainMap)
	bytes := buffer.Bytes()
	file.Write(bytes)
	zap.L().Info("Snapshot taken successfully")
	file.Close()
}
func RunSnapShotTaker(mainMap schemas.MainMap) {
	var wg sync.WaitGroup
	wg.Add(1)
	go takeSnapShot(&wg, mainMap)
	wg.Wait()
}
