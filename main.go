package main

import (
	"fmt"
	"in-memory-store/schemas"
	"in-memory-store/snapshots"
)

func main() {
	logger := GetLogger()
	defer logger.Sync()
	globalMap := schemas.CreateMainMap()
	logger.Info("Application Initilized")
	snapshots.ReadSnapShotFile(&globalMap)
	fmt.Println("Main map is", globalMap)

	// globalMap.SetInteger("thisisveryveryveryveryveryverylong", int64(458234092380598235))
	// globalMap.SetInteger("minusOne", int64(-1))
	// globalMap.SetInteger("plus2", int64(2))
	// globalMap.SetInteger("plus3", int64(3))
	// globalMap.SetInteger("plus10", int64(10))
	// globalMap.SetInteger("one hundred", int64(100))
	// globalMap.SetInteger("hehehehe", int64(69))
	// globalMap.SetString("firstString", "this is the value of the first string")
	// globalMap.SetString("New string", "new string")
	arr := []int64{10, 20, 30, 40, 50}
	globalMap.SetIntegerArray("int arr", arr)
	arr2 := []int64{80, 90, 500, 200, 80808080}
	globalMap.SetIntegerArray("int arr12", arr2)
	snapshots.RunSnapShotTaker(globalMap)
	logger.Info("Application Closing....")
}
