package orm

import (
	"bytes"
	"github.com/xwb1989/sqlparser"
	"io"
	"os"
)

func (G *Gorm) Parser(file string) ([]byte, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return b, err
	}
	tokens := sqlparser.NewTokenizer(bytes.NewReader(b))
	for {
		_, err := sqlparser.ParseNext(tokens)

		if err == io.EOF {
			return b, nil
		}
		if err != nil {
			return b, err
		}
	}
}
