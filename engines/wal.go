package engines

import (
	"encoding/base64"
	"fmt"
	"bytes"
	"io"
	"os"
	"path"
	"encoding/gob"
)

// clearWAL closes the current file and open the new file in the truncate mode.
func clearWAL(dbDir string, wal *os.File) (*os.File, error) {
	walPath := path.Join(dbDir, walFileName)

	if err := wal.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the WAL file %s: %w", walPath, err)
	}

	wal, err := os.OpenFile(walPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file %s: %w", walPath, err)
	}

	return wal, nil
}


// appendToWAL appends entry to the WAL file.
func appendToWAL(wal *os.File, row Row) error {
	// for safety, since the file is open in read-write mode
	gob.Register(Row{})
	if _, err := wal.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to the end: %w", err)
	}

	if err := encode(row, wal); err != nil {
		return fmt.Errorf("failed to encode and write to the file: %w", err)
	}

	if err := wal.Sync(); err != nil {
		return fmt.Errorf("failed to sync the file: %w", err)
	}

	return nil
}

// loadMemTable loads MemTable from the WAL file.
func loadMemTable(wal *os.File) (*memTable, error) {
	// for safety, since the file is open in read-write mode
	if _, err := wal.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to the beginning: %w", err)
	}

	memTable := newMemTable()
	for {
		row, err := decode(wal)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read: %w", err)
		}
		if err == io.EOF {
			return memTable, nil
		}

		if value != nil {
			memTable.put(row)
		} else {
			memTable.delete(row)
		}
	}
}


func ToGOB64(m Row) string {
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(m)
    if err != nil { fmt.Println(`failed gob Encode`, err) }
    return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) Row {
    m := Row{}
    by, err := base64.StdEncoding.DecodeString(str)
    if err != nil { fmt.Println(`failed base64 Decode`, err); }
    b := bytes.Buffer{}
    b.Write(by)
    d := gob.NewDecoder(&b)
    err = d.Decode(&m)
    if err != nil { fmt.Println(`failed gob Decode`, err); }
    return m
}

func encode(row Row, wal *os.File) error {
	gob_row := ToGOB64(row)
	_, err := wal.Write([]byte(gob_row))
	if err != nil {
		return fmt.Errorf("failed to write gob to the file: %w", err)
	}
	return nil
}


func decode(wal *os.File) (Row, error) {
	gob_row := make([]byte, 0)
	_, err := wal.Read(gob_row)
	if err != nil {
		return Row{}, err
	}
	row := FromGOB64(string(gob_row))
	return row, nil
}