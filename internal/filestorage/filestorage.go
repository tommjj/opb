package filestorage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"reflect"
	"sync"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrTypeNotPtr      = errors.New("type must a pointer")
	ErrNotExist        = os.ErrNotExist
)

type FileStorage struct {
	Filename string
	lock     sync.RWMutex
}

func New(filename string) *FileStorage {
	return &FileStorage{
		Filename: filename,
	}
}

func (f *FileStorage) Load(v any) error {
	f.lock.RLock()
	defer f.lock.RUnlock()

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrTypeNotPtr
	}

	kind := val.Elem().Kind()
	if kind != reflect.Slice && kind != reflect.Struct && kind != reflect.Map {
		return ErrUnsupportedType
	}

	file, err := os.Open(f.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) Sync(v any) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(f.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
