package MineDb

import (
	"encoding/json"
	"fmt"
	"os"
)

var file_db_name string = "mine_db.json"

type DBError string

func (e DBError) Error() string {
	return string(e)
}

func set_db(mine_db map[string]any) error {
	to_write, err := json.Marshal(mine_db)
	if err != nil {
		return err
	}
	err = os.WriteFile(file_db_name, to_write, 0644)
	return err
}

func get_db() (map[string]any, error) {
	var ret_map map[string]any
	ret_map = make(map[string]any)
	res, err := os.Open(file_db_name)
	if err == nil {
		cur_dec := json.NewDecoder(res)
		cur_dec.Decode(&ret_map)
		res.Close()
	} else {
		set_db(ret_map)
		return get_db()
	}
	return ret_map, err
}

func Get(key string) (any, error) {
	mine_db, err := get_db()
	if err != nil {
		return nil, err
	}
	ans, ok := mine_db[key]
	if !ok {
		return nil, DBError(fmt.Sprintf("the following key is not in the mine_db: %s", key))
	}
	return ans, nil
}

func Set(key string, v any) error {
	mine_db, err := get_db()
	if err != nil {
		return err
	}
	mine_db[key] = v
	err = set_db(mine_db)
	return err
}
