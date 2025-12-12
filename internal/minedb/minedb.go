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

func (curDb *mineDb) saveDb() error {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	toWrite, err := json.Marshal(curDb.dbMap)
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
		err = curDb.saveDb()
		curDb.mu.Lock()
		return retMap, err
	}
	return retMap, err
}

func (curDb *mineDb) getVal(key string) (any, error) {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	ans, ok := curDb.dbMap[key]
	if !ok {
		return nil, DBError(fmt.Sprintf("the following key is not in the mineDb: %s", key))
	}
	return ans, nil
}

func (curDb *mineDb) setKV(key string, v any) {
	curDb.mu.Lock()
	defer curDb.mu.Unlock()
	curDb.dbMap[key] = v
}

func (curDb *mineDb) Get(key string, v *any) error {
	ans, err := curDb.getVal(key)
	if err != nil {
		return err
	}
	MarshalledV, err := json.Marshal(ans)
	if err != nil {
		return err
	}
	err = json.Unmarshal(MarshalledV, v)
	if err != nil {
		return err
	}
	*v = ans // somehow this marshal unmarshall strategy still sometimes doesnt convert the type to the needed one, so I decided to only check for errors with it and stick with this method for actual assignment
	return err
}

func (curDb *mineDb) Set(key string, v any) error {
	curDb.setKV(key, v)
	err := curDb.saveDb()
	return err
}

var once sync.Once
var curMineDb *mineDb

func NewMineDb() (*mineDb, error) {
	var err error = nil
	once.Do(func() {
		curMineDb = &mineDb{}
		curMineDb.dbMap, err = curMineDb.getDb()
	})
	return curMineDb, err
}
