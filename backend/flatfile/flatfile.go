package flatfile

import (
    "encoding/csv"
    "os"
)

type FlatFile struct {
    Path      string
    Delimiter string
}

func (f *FlatFile) GetColumns() ([]string, error) {
    file, err := os.Open(f.Path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = rune(f.Delimiter[0])
    headers, err := reader.Read()
    if err != nil {
        return nil, err
    }
    return headers, nil
}