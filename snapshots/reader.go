package snapshots

import (
	"encoding/binary"
	"fmt"
	"in-memory-store/constants"
	"in-memory-store/schemas"
	"io"
	"os"

	"go.uber.org/zap"
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
	bytes := make([]byte, len(constants.FILE_HEADER))
	l, err := reader.file.Read(bytes)
	if err != nil {
		return err
	}
	if l < len(constants.FILE_HEADER) {
		return fmt.Errorf("Error while skipping file header expected %d bytes found %d", len(constants.FILE_HEADER), l)
	}
	header_contents := string(bytes)
	if header_contents != constants.FILE_HEADER {
		return fmt.Errorf("Invalid file header found, expected %s, found %s", constants.FILE_HEADER, header_contents)
	}
	return err
}

func (reader *BinaryReader) skipBlockSeperator() error {
	bytes := make([]byte, constants.BLOCK_SEPERATOR_LENGTH)
	l, err := reader.file.Read(bytes)
	if err != nil {
		return err
	}
	if l < constants.BLOCK_SEPERATOR_LENGTH {
		return fmt.Errorf("Error skipping block seperator not enough bytes found")
	}
	blockSeperatorString := string(bytes)
	if blockSeperatorString != "\r\n" {
		return fmt.Errorf("Block seperator expected '\r\n' found %s", blockSeperatorString)
	}
	return nil
}

func (reader *BinaryReader) getInt64DataFromBlock() (int64, error) {
	bytes := make([]byte, constants.INT_TYPE_LENGTH)
	l, err := reader.file.Read(bytes)
	if err != nil {
		return 0, err
	}
	if l < constants.INT_TYPE_LENGTH {
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
func handleError(err error, context string) bool {
	if err == io.EOF {
		return true
	}
	if err != nil {
		zap.L().Error(context, zap.Error(err))
		return true
	}
	return false
}

func ReadSnapShotFile(mainMap *schemas.MainMap) {
	f, err := os.Open(constants.SNAPSHOT_FILE_NAME)
	defer f.Close()
	if err != nil {
		zap.L().Error("Error reading snapshot file", zap.Error(err))
		return
	}
	reader := CreateBinaryReader(f)
	err = reader.skipFileHeader()
	if err != nil {
		zap.L().Error("Error skipping file header", zap.Error(err))
	}
	reader.skipBlockSeperator()
	for {
		blockValueType, err := reader.getInt64DataFromBlock()
		if handleError(err, "Error while reading block type") {
			return
		}
		keyLength, err := reader.getInt64DataFromBlock()
		if handleError(err, "Error while reading key length") {
			return
		}

		key, err := reader.getStringDataFromBlock(keyLength)
		if handleError(err, "Error while reading key") {
			return
		}

		if blockValueType == int64(constants.INTEGER_TYPE) {
			blockValue, err := reader.getInt64DataFromBlock()
			if handleError(err, "Error while reading block value") {
				return
			}
			mainMap.SetInteger(key, blockValue)
		} else if blockValueType == int64(constants.STRING_TYPE) {
			valueLength, err := reader.getInt64DataFromBlock()
			if handleError(err, "Error while reading block value") {
				return
			}
			blockValue, err := reader.getStringDataFromBlock(valueLength)
			if handleError(err, "Error while reading block value") {
				return
			}
			mainMap.SetString(key, blockValue)
		}

		err = reader.skipBlockSeperator()
		if handleError(err, "Error while skipping block") {
			return
		}
	}
}
