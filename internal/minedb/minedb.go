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
	mu sync.Mutex // add for ALL fucntions
}

func (curDb *mineDb) setDb(stateToWrite map[string]any) error {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	toWrite, err := json.Marshal(stateToWrite)
	if err != nil {
		return err
	}
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

func (curDb *mineDb) get(key string) (any, error) {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	curDb.mu.Unlock()
	curState, err := curDb.getDb()
	curDb.mu.Lock()
	if err != nil {
		return nil, err
	}
	ans, ok := curState[key]
	if !ok {
		return nil, DBError(fmt.Sprintf("the following key is not in the mineDb: %s", key))
	}
	return ans, nil
}

func (curDb *mineDb) set(key string, v any) error {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	curDb.mu.Unlock()
	mineDb, err := curDb.getDb()
	curDb.mu.Lock()
	if err != nil {
		return err
	}
	mineDb[key] = v
	curDb.mu.Unlock()
	err = curDb.setDb(mineDb)
	curDb.mu.Lock()
	return err
}

var once sync.Once
var curMineDb *mineDb

func getMineDb() *mineDb {
	once.Do(func() {
		curMineDb = &mineDb{}
	})
	return curMineDb
}

func Set(key string, v any) error {
	curLink := getMineDb()
	return curLink.set(key, v)
}

func Get(key string) (any, error) {
	curLink := getMineDb()
	return curLink.get(key)
}
