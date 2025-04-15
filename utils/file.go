package utils

import (
	"bytes"
	"encoding/csv"
)

const MaxCreateDraftSize = 1024 * 1024

func GenerateFileData(table [][]string) ([]byte, error) {

	buf := &bytes.Buffer{}
	wr := csv.NewWriter(buf)

	err := wr.WriteAll(table)
	if err != nil {
		return []byte{}, err
	}

	defer wr.Flush()

	return buf.Bytes(), nil
}
