package main

import (
	"go.uber.org/zap"
)

func CreateMainMap() MainMap {
	return MainMap{
		INTEGER_MAP:       make(map[string]int64),
		STRING_MAP:        make(map[string]string),
		STRING_ARRAY_MAP:  make(map[string][]string),
		INTEGER_ARRAY_MAP: make(map[string][]int64),
	}
}

func (m *MainMap) setInteger(key string, value int64) {
	zap.L().Info("Setting Integer", zap.String("key", key), zap.Int64("value", value))
	m.INTEGER_MAP[key] = value
}

func (m *MainMap) setString(key string, value string) {
	zap.L().Info("Setting String", zap.String("key", key), zap.String("value", value))
	m.STRING_MAP[key] = value
}

func (m *MainMap) setIntegerArray(key string, value []int64) {
	zap.L().Info("Setting Integer Array", zap.String("key", key), zap.Int64s("value", value))
	m.INTEGER_ARRAY_MAP[key] = value
}

func (m *MainMap) setStringArray(key string, value []string) {
	zap.L().Info("Setting String Array", zap.String("key", key), zap.Strings("value", value))
	m.STRING_ARRAY_MAP[key] = value
}

func (m *MainMap) getInteger(key string) int64 {
	return m.INTEGER_MAP[key]
}

func (m *MainMap) getString(key string) string {
	return m.STRING_MAP[key]
}
func (m *MainMap) getIntegerArray(key string) []int64 {
	return m.INTEGER_ARRAY_MAP[key]
}

func (m *MainMap) getStringArray(key string) []string {
	return m.STRING_ARRAY_MAP[key]
}

func (m *MainMap) getValue(key string) interface{} {
	stringValue := m.getString(key)
	if stringValue != "" {
		return stringValue
	}
	intValue := m.getInteger(key)
	if intValue != 0 {
		return intValue
	}
	stringArrayValue := m.getStringArray(key)
	if len(stringArrayValue) > 0 {
		return intValue
	}
	intArrayValue := m.getIntegerArray(key)
	if len(stringArrayValue) > 0 {
		return intArrayValue
	}
	return nil

}
