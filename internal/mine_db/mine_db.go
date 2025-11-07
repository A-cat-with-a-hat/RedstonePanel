package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var file_db_name string = "mine_db.json"

type MyError struct {
	What string
	When time.Time
}

func (e MyError) Error() string {
	return fmt.Sprintf("Error %s, at %v", e.What, e.When)
}

func set_db(mine_db map[string]any) error {
	to_write, cur_error := json.Marshal(mine_db)
	if cur_error != nil {
		return cur_error
	}
	cur_error = os.WriteFile(file_db_name, to_write, 0644)
	// I decided to go with 0644 instead of 0666, because this file shouldn't be modified outside the programm
	// I still wanted others to be able to see its content so I didn't go with 0600
	return cur_error
}

func get_db() (map[string]any, error) {
	var ret_map map[string]any
	ret_map = make(map[string]any)
	res, cur_error := os.Open(file_db_name)
	if cur_error == nil {
		cur_dec := json.NewDecoder(res)
		cur_dec.Decode(&ret_map)
		res.Close()
	} else {
		set_db(ret_map)
		return get_db()
	}
	return ret_map, cur_error
}

func Get(key string) (any, error) {
	mine_db, cur_error := get_db()
	if cur_error != nil {
		return nil, cur_error
	}
	ans, ok := mine_db[key]
	if !ok {
		return nil, MyError{fmt.Sprintf("the following key is not in the mine_db: %s", key), time.Now()}
	}
	return ans, nil
}

func Set(key string, v any) error {
	mine_db, cur_error := get_db()
	if cur_error != nil {
		return cur_error
	}
	mine_db[key] = v
	cur_error = set_db(mine_db)
	return cur_error
}
