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

	globalMap.setInteger("thisisveryveryveryveryveryverylong", int64(458234092380598235))
	globalMap.setInteger("minusOne", int64(-1))
	globalMap.setInteger("plus2", int64(2))
	globalMap.setInteger("plus3", int64(3))
	// globalMap.setInteger("plus10", int64(10))
	// globalMap.setInteger("one hundred", int64(100))
	// globalMap.setInteger("hehehehe", int64(69))

	RunSnapShotTaker(globalMap)
	logger.Info("Application Closing....")
}
