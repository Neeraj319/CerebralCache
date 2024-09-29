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
	// globalMap.setStringArray("names", []string{"hello", "how are you"})
	// globalMap.setStringArray("songs", []string{"wish you were here", "as tears go by"})

	globalMap.setInteger("thisisveryveryveryveryveryverylong", int64(458234092380598235))
	globalMap.setInteger("minusOne", int64(-1))
	globalMap.setInteger("plus2", int64(2))
	globalMap.setInteger("plus3", int64(3))
	// globalMap.setString("name", "hero")

	RunSnapShotTaker(globalMap)
	logger.Info("Application Closing....")
}
