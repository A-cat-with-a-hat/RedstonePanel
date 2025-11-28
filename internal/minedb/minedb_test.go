package minedb

import (
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func checkGet(t *testing.T, curDb *mineDb, key string, expectedValue any) {
	var valGet = expectedValue
	err := curDb.Get(key, &valGet)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	if !reflect.DeepEqual(valGet, expectedValue) {
		t.Errorf("With key = %s value %v is needed, but got: %v\n", key, expectedValue, valGet)
	}
}

func checkSet(t *testing.T, curDb *mineDb, key string, value any) {
	err := curDb.Set(key, value)
	if err != nil {
		t.Errorf("%s\n", err)
	}
}

type concurrencySafeMap struct {
	mu     sync.Mutex
	curMap map[string]any
}

func (curSafeMap *concurrencySafeMap) setVal(key string, val any) {
	curSafeMap.mu.Lock()
	curSafeMap.curMap[key] = val
	curSafeMap.mu.Unlock()
}

func TestMineDb(t *testing.T) {
	curLink, err := GetMineDb()
	if err != nil {
		t.Errorf("%s\n", err)
	}
	keys := []string{"Time", "Location", "Mode", "What", "Doctor", "Ostrich"}
	values := []any{float64(123.34), 38, "who", [3]int{1, 2, 3}, 342, "AAAA", nil} // these should probably be random but i dont care
	curRandSource := rand.NewSource(38)
	randFunc := rand.New(curRandSource)
	var testSize int = 100
	var correctMineDb = &(concurrencySafeMap{})
	correctMineDb.curMap = make(map[string]any)
	for i := 0; i < testSize; i++ {
		go func() { // test without a grace period for runtime errors
			curKey := keys[randFunc.Intn(len(keys))]
			curVal := values[randFunc.Intn(len(values))]
			checkSet(t, curLink, curKey, curVal)
		}()
	}
	for i := 0; i < testSize; i++ { // test without multiflow
		curKey := keys[randFunc.Intn(len(keys))]
		curVal := values[randFunc.Intn(len(values))]
		correctMineDb.setVal(curKey, curVal)
		checkSet(t, curLink, curKey, curVal)
	}
	for k := range correctMineDb.curMap {
		var valToGet any = correctMineDb.curMap[k]
		checkGet(t, curLink, k, valToGet)
	}

	for i := 0; i < testSize; i++ {
		go func() {
			curKey := keys[randFunc.Intn(len(keys))]
			curVal := values[randFunc.Intn(len(values))]
			correctMineDb.setVal(curKey, curVal)
			checkSet(t, curLink, curKey, curVal)
		}()
		time.Sleep(50 * time.Millisecond) // test wtih a grace period to check that mineDb works fine multiflow. Is it needed? I dont know, but without it the tests fail)
	}
	for k := range correctMineDb.curMap {
		var valToGet any = correctMineDb.curMap[k]
		checkGet(t, curLink, k, valToGet)
	}
}
