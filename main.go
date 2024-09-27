package main

import (
	"go.uber.org/zap"
)

const (
	INTEGER_MAP = iota
	STRING_MAP
	INTEGER_ARRAY_MAP
	STRING_ARRAY_MAP
)

type MainMap struct {
	INTEGER_MAP       map[string]int
	STRING_MAP        map[string]string
	INTEGER_ARRAY_MAP map[string][]int
	STRING_ARRAY_MAP  map[string][]string
}

func createMainMap() MainMap {
	return MainMap{
		INTEGER_MAP:       make(map[string]int),
		STRING_MAP:        make(map[string]string),
		STRING_ARRAY_MAP:  make(map[string][]string),
		INTEGER_ARRAY_MAP: make(map[string][]int),
	}
}

func (m *MainMap) setInteger(key string, value int) {
	zap.L().Info("Setting Integer", zap.String("key", key), zap.Int("value", value))
	m.INTEGER_MAP[key] = value
}

func (m *MainMap) setString(key string, value string) {
	zap.L().Info("Setting String", zap.String("key", key), zap.String("value", value))
	m.STRING_MAP[key] = value
}

func (m *MainMap) setIntegerArray(key string, value []int) {
	zap.L().Info("Setting Integer Array", zap.String("key", key), zap.Ints("value", value))
	m.INTEGER_ARRAY_MAP[key] = value
}

func (m *MainMap) setStringArray(key string, value []string) {
	zap.L().Info("Setting String Array", zap.String("key", key), zap.Strings("value", value))
	m.STRING_ARRAY_MAP[key] = value
}

func main() {
	logger := GetLogger()
	defer logger.Sync()
	global_map := createMainMap()
	logger.Info("Application Initilized")
	global_map.setStringArray("names", []string{"hello", "how are you"})
	global_map.setStringArray("songs", []string{"wish you were here", "as tears go by"})

	logger.Info("Application Closing....")
}
