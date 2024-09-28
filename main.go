package main

type MainMap struct {
	INTEGER_MAP       map[string]int64
	STRING_MAP        map[string]string
	INTEGER_ARRAY_MAP map[string][]int64
	STRING_ARRAY_MAP  map[string][]string
}

func main() {
	logger := GetLogger()
	defer logger.Sync()
	globalMap := CreateMainMap()
	logger.Info("Application Initilized")
	globalMap.setStringArray("names", []string{"hello", "how are you"})
	globalMap.setStringArray("songs", []string{"wish you were here", "as tears go by"})
	globalMap.setInteger("hehehe", 1)
	globalMap.setString("name", "hero")
	globalMap.setInteger("minusOne", -1)
	globalMap.setInteger("plus2", 2)
	globalMap.setInteger("plus3", 3)

	RunSnapShotTaker(globalMap)
	logger.Info("Application Closing....")
}
