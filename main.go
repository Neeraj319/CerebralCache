package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type MainMap struct {
	INTEGER_MAP       map[string]int64
	STRING_MAP        map[string]string
	INTEGER_ARRAY_MAP map[string][]int64
	STRING_ARRAY_MAP  map[string][]string
}

func main() {
	logger := GetLogger()
	f, _ := os.Create("file.bin")
	defer logger.Sync()
	global_map := CreateMainMap()
	logger.Info("Application Initilized")
	global_map.setStringArray("names", []string{"hello", "how are you"})
	global_map.setStringArray("songs", []string{"wish you were here", "as tears go by"})
	global_map.setInteger("hehehe", 1)
	global_map.setString("name", "hero")

	logger.Info("Application Closing....")
	a := []byte(global_map.getString("name"))
	fmt.Println("a is", a)
	err := binary.Write(f, binary.LittleEndian, global_map.getString("name"))
	fmt.Println("error is", err)
	f.Close()
}
