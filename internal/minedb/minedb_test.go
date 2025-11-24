package minedb

import (
	"math/rand"
	"reflect"
	"testing"
)

func checkGet(t *testing.T, key string, expectedValue any) {
	valGet, err := Get(key)
	if err != nil {
		t.Errorf("%s", err)
	}
	if !reflect.DeepEqual(valGet, expectedValue) {
		t.Errorf("Need value %v, got: %v", expectedValue, valGet)
	}
}

func checkSet(t *testing.T, key string, value any) {
	err := Set(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	checkGet(t, key, value)
}

func TestMineDb(t *testing.T) {
	keys := []string{"Time", "Location", "Mode", "What", "Doctor", "Ostrich"}
	values := []any{"2 am", float64(38), "who", []any{float64(1), float64(2), float64(3)}, float64(342), "AAAA", nil} // they should probably be random but i just dont care
	// Turns out that json turns all numbers into float64 and all slices into []any see here for more information: https://pkg.go.dev/encoding/json#Unmarshal
	// Why this is true is beyond me and I cannot fix this problem, so I hope everyone likes working with floats and doesnt care what they're storing in their arrays
	curRandSource := rand.NewSource(38)
	randFunc := rand.New(curRandSource)
	var correctMineDb map[string]any = make(map[string]any)
	var testSize int = 100
	usedKeys := []string{}
	for i := 0; i < testSize; i++ {
		curKey := keys[randFunc.Intn(len(keys))]
		curVal := values[randFunc.Intn(len(values))]
		correctMineDb[curKey] = curVal
		usedKeys = append(usedKeys, curKey)
		checkSet(t, curKey, curVal)
	}
	for i := 0; i < testSize; i++ {
		checkGet(t, usedKeys[i], correctMineDb[usedKeys[i]])
	}
}
