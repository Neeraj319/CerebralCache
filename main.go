package main

import "fmt"

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
	m.INTEGER_MAP[key] = value
}

func (m *MainMap) setString(key string, value string) {
	m.STRING_MAP[key] = value
}

func (m *MainMap) setIntegerArray(key string, value []int) {
	m.INTEGER_ARRAY_MAP[key] = value
}

func (m *MainMap) setStringArray(key string, value []string) {
	m.STRING_ARRAY_MAP[key] = value
}

func main() {
	global_map := createMainMap()
	global_map.setStringArray("names", []string{"hello", "how are you"})
	global_map.setStringArray("songs", []string{"wish you were here", "as tears go by"})
	fmt.Println("string array map", global_map.STRING_ARRAY_MAP)
}
