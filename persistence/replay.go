package persistence

import (
	"errors"
	"fmt"
	"github.com/krishsinghhura/goredis/resp"
	"io"
	"os"
)

func ReplayAOF(path string, apply func(args []string) error) error {
	return replay(path, false, apply)
}

func ReplayAOFTruncateTail(path string, apply func(args []string) error) error {
	return replay(path, true, apply)
}

func replay(path string, truncateTail bool, apply func(args []string) error) error {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	p := resp.NewParser(f)
	var lastGoodOffset int64

	for {

		cur, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		lastGoodOffset = cur

		args, err := readOneRESPCommand(p)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			if truncateTail && errors.Is(err, io.ErrUnexpectedEOF) {
				_ = f.Close()
				return os.Truncate(path, lastGoodOffset)
			}

			return fmt.Errorf("AOF parse error at offset %d: %w", lastGoodOffset, err)
		}

		if len(args) == 0 {
			continue
		}

		if err := apply(args); err != nil {
			return fmt.Errorf("AOF apply error at offset %d: %w", lastGoodOffset, err)
		}
	}
}

func readOneRESPCommand(p *resp.Parser) ([]string, error) {
	v, err := p.ParseOne()
	if err != nil {
		return nil, err
	}

	if v.Typ != "Array" {
		return nil, fmt.Errorf("invalid AOF entry: expected Array, got %s", v.Typ)
	}

	args := make([]string, 0, len(v.Array))
	for i, elem := range v.Array {
		if elem.Typ != "BulkString" {
			return nil, fmt.Errorf("invalid AOF entry: arg %d is %s, expected BulkString", i, elem.Typ)
		}
		args = append(args, elem.Str)
	}

	return args, nil
}
