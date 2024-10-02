package schemas

import "go.uber.org/zap"

type MainMap struct {
	INTEGER_MAP       map[string]int64
	STRING_MAP        map[string]string
	INTEGER_ARRAY_MAP map[string][]int64
	STRING_ARRAY_MAP  map[string][]string
}

func CreateMainMap() MainMap {
	return MainMap{
		INTEGER_MAP:       make(map[string]int64),
		STRING_MAP:        make(map[string]string),
		INTEGER_ARRAY_MAP: make(map[string][]int64),
		STRING_ARRAY_MAP:  make(map[string][]string),
	}
}

func (m *MainMap) SetInteger(key string, value int64) {
	zap.L().Info("Setting Integer", zap.String("key", key), zap.Int64("value", value))
	m.INTEGER_MAP[key] = value
}

func (m *MainMap) SetString(key string, value string) {
	zap.L().Info("Setting String", zap.String("key", key), zap.String("value", value))
	m.STRING_MAP[key] = value
}

func (m *MainMap) SetIntegerArray(key string, value []int64) {
	zap.L().Info("Setting Integer Array", zap.String("key", key), zap.Int64s("value", value))
	m.INTEGER_ARRAY_MAP[key] = value
}

func (m *MainMap) SetStringArray(key string, value []string) {
	zap.L().Info("Setting String Array", zap.String("key", key), zap.Strings("value", value))
	m.STRING_ARRAY_MAP[key] = value
}

func (m *MainMap) GetInteger(key string) int64 {
	return m.INTEGER_MAP[key]
}

func (m *MainMap) GetString(key string) string {
	return m.STRING_MAP[key]
}
func (m *MainMap) GetIntegerArray(key string) []int64 {
	return m.INTEGER_ARRAY_MAP[key]
}

func (m *MainMap) getStringArray(key string) []string {
	return m.STRING_ARRAY_MAP[key]
}

func (m *MainMap) getValue(key string) interface{} {
	stringValue := m.GetString(key)
	if stringValue != "" {
		return stringValue
	}
	intValue := m.GetInteger(key)
	if intValue != 0 {
		return intValue
	}
	stringArrayValue := m.getStringArray(key)
	if len(stringArrayValue) > 0 {
		return intValue
	}
	intArrayValue := m.GetIntegerArray(key)
	if len(stringArrayValue) > 0 {
		return intArrayValue
	}
	return nil

}
