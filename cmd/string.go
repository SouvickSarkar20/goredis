package cmd

import (
	"github.com/SouvickSarkar20/goredis/resp"
	"github.com/SouvickSarkar20/goredis/store"
	"strconv"
	"strings"
	"time"
)

func handleSet(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) < 3 {
		return w.WriteError("ERR : WRONG NO OF ARGUMENTS")
	}

	key := input.Array[1].Str
	value := input.Array[2].Str

	var duration time.Duration

	if len(input.Array) >= 5 {
		option := strings.ToUpper(input.Array[3].Str)
		if option == "EX" {
			seconds, err := strconv.Atoi(input.Array[4].Str)
			if err != nil {
				return w.WriteError("ERR : value is not an integer")
			}
			duration = time.Duration(seconds) * time.Second
		}
	}

	db.Set(key, value, duration)
	return w.WriteSimpleString("OK")
}

func handleGet(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) != 2 {
		return w.WriteError("ERR : wrong no of arguments")
	}

	key := input.Array[1].Str
	value, exists := db.Get(key)

	if !exists {
		return w.WriteBulkString("") // Return empty string if key does not exist
	}

	strVal, ok := value.(string)
	if !ok {
		return w.WriteError("ERR : The value could not be converted to string")
	}

	return w.WriteBulkString(strVal)
}
