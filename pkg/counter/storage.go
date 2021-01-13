package counter

import (
	"container/list"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

type FileStorage struct {
	fileName string
}

func NewFileStorage(fileName string) *FileStorage {
	return &FileStorage{
		fileName: fileName,
	}
}

func (f *FileStorage) Save(events []time.Time) error {
	file, err := os.Create(f.fileName)
	if err != nil {
		return errors.Wrap(err, "can not open file")
	}
	defer file.Close()

	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		if _, err := file.Write([]byte(event.Format(time.RFC1123Z))); err != nil {
			return errors.Wrap(err, "can not write line")
		}
	}

	return nil
}

func (f *FileStorage) Get() (*list.List, error) {
	events := list.New()
	file, err := os.Open(f.fileName)
	if err != nil {
		return events, errors.Wrap(err, "can not write file")
	}
	defer file.Close()

	dtRaw := make([]byte, len(time.RFC1123Z))
	for {
		if _, err := file.Read(dtRaw); err != nil {
			if err == io.EOF {
				break
			}
			return events, err
		}
		dt, err := time.Parse(time.RFC1123Z, string(dtRaw))
		if err != nil {
			return events, err
		}
		events.PushBack(dt)
	}

	return events, nil
}
