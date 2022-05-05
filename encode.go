package jse

import (
	"bytes"
	"encoding/csv"
)

func (e Event) Encode() string {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	_ = w.Write(e.record())
	w.Flush()
	return b.String()
}
