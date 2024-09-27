package main

type MainMap struct {
	INTEGER_MAP       map[string]int
	STRING_MAP        map[string]string
	INTEGER_ARRAY_MAP map[string][]int
	STRING_ARRAY_MAP  map[string][]string
}

func main() {
	logger := GetLogger()
	defer logger.Sync()
	global_map := CreateMainMap()
	logger.Info("Application Initilized")
	global_map.setStringArray("names", []string{"hello", "how are you"})
	global_map.setStringArray("songs", []string{"wish you were here", "as tears go by"})
	global_map.setInteger("hehehe", 1)

	logger.Info("Application Closing....")
}
