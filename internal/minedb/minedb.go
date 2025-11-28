package minedb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var fileDbName string = "minedb.json"

type DBError string

func (e DBError) Error() string {
	return string(e)
}

type mineDb struct {
	mu    sync.Mutex
	dbMap map[string]any
}

func (curDb *mineDb) setDb(stateToWrite map[string]any) error {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	toWrite, err := json.Marshal(stateToWrite)
	//	fmt.Printf("Error Set = %s\n", err)
	if err != nil {
		return err
	}
	curDb.dbMap = stateToWrite
	err = os.WriteFile(fileDbName, toWrite, 0644)
	return err
}

func (curDb *mineDb) getDb() (map[string]any, error) {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	var retMap map[string]any = make(map[string]any)
	res, err := os.Open(fileDbName)
	if err == nil {
		curDec := json.NewDecoder(res)
		curDec.Decode(&retMap)
		res.Close()
	} else {
		curDb.mu.Unlock()
		err = curDb.setDb(retMap)
		curDb.mu.Lock()
		return retMap, err
	}
	return retMap, err
}

func (curDb *mineDb) Get(key string, v *any) error {
	ans, ok := curDb.dbMap[key]
	*v = ans // i'm not sure how to handle errors with this assignment and still be able to convert types when possible(float32 to float64 or something like that)
	if !ok {
		return DBError(fmt.Sprintf("the following key is not in the mineDb: %s", key))
	}
	return nil
}

func (curDb *mineDb) Set(key string, v any) error {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	curDb.dbMap[key] = v
	curDb.mu.Unlock()
	err := curDb.setDb(curDb.dbMap)
	curDb.mu.Lock()
	return err
}

var once sync.Once
var curMineDb *mineDb

func getMineDb() (*mineDb, error) {
	var err error = nil
	once.Do(func() {
		curMineDb = &mineDb{}
		curMineDb.dbMap, err = curMineDb.getDb()
	})
	return curMineDb, err
}
