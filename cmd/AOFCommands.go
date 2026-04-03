package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SouvickSarkar20/goredis/store"
)

// ApplyAOFCommand applies one replayed command directly to store.
func ApplyAOFCommand(db *store.Store, args []string) error {
	if len(args) == 0 {
		return nil
	}

	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		return nil // no state change

	case "SET":
		// support: SET key value
		// support: SET key value EX seconds
		if len(args) != 3 && len(args) != 5 {
			return fmt.Errorf("invalid SET args in AOF")
		}

		key := args[1]
		value := args[2]

		if len(args) == 3 {
			db.Set(key, value, 0)
			return nil
		}

		opt := strings.ToUpper(args[3])
		raw := args[4]

		if opt != "EX" {
			return fmt.Errorf("unsupported SET option in AOF: %s", opt)
		}
		sec, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("invalid EX value: %w", err)
		}
		db.Set(key, value, time.Duration(sec)*time.Second)
		return nil

	case "DEL":
		if len(args) != 2 {
			return fmt.Errorf("invalid DEL args in AOF")
		}
		key := args[1]
		db.Delete(key)
		return nil

	case "LPUSH":
		if len(args) != 3 {
			return fmt.Errorf("invalid LPUSH args in AOF")
		}
		key := args[1]
		value := args[2]
		return db.LPush(key, value)

	case "LPOP":
		if len(args) != 2 {
			return fmt.Errorf("invalid LPOP args in AOF")
		}
		key := args[1]
		db.LPop(key)
		return nil

	case "HSET":
		if len(args) != 4 {
			return fmt.Errorf("invalid HSET args in AOF")
		}
		return db.HSet(args[1], args[2], args[3])

	case "HDEL":
		if len(args) != 3 {
			return fmt.Errorf("invalid HDEL args in AOF")
		}
		_, err := db.HDel(args[1], args[2])
		return err

	case "SADD":
		if len(args) != 3 {
			return fmt.Errorf("invalid SADD args in AOF")
		}
		return db.SAdd(args[1], args[2])

	case "SREM":
		if len(args) != 3 {
			return fmt.Errorf("invalid SREM args in AOF")
		}
		_, err := db.SRem(args[1], args[2])
		return err

	// read-only commands in AOF can be ignored safely
	case "GET", "SMEMBERS", "SISMEMBER", "HGET":
		return nil

	default:
		return fmt.Errorf("unsupported command in AOF: %s", cmd)
	}
}
