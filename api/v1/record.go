package v1

import (
	"encoding/binary"
	"io"
)

type Record struct {
	Value  []byte
	Offset uint64
}

func (r *Record) Marshal(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, r.Offset); err != nil {
		return err
	}
	if _, err := w.Write(r.Value); err != nil {
		return err
	}
	return nil
}

func (r *Record) Unmarshal(rd io.Reader) error {
	if err := binary.Read(rd, binary.BigEndian, &r.Offset); err != nil {
		return err
	}
	var err error
	r.Value, err = io.ReadAll(rd)
	if err != nil {
		return err
	}
	return nil
}
