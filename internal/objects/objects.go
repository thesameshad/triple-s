package objects

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type Object struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified string
}

var (
	objectKeyRe   = regexp.MustCompile(`^[a-zA-Z0-9!_\-\.]\S{0,255}$`)
	ErrKeyInvalid = errors.New("Invalid object key")
	ErrNotFound   = errors.New("Object not found")
)

func isValidObjectKey(key string) bool {
	if len(key) < 1 || len(key) > 255 {
		return false
	}
	return true
}

func ObjectsFilePath(dataDir, bucket string) string {
	return filepath.Join(dataDir, bucket, "objects.csv")
}

func LoadObjects(dataDir, bucket string) ([]Object, error) {
	f, err := os.Open(ObjectsFilePath(dataDir, bucket))
	if os.IsNotExist(err) {
		return []Object{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var objs []Object
	for _, rec := range records {
		if len(rec) < 4 {
			continue
		}
		sz, _ := strconv.ParseInt(rec[1], 10, 64)
		objs = append(objs, Object{
			Key:          rec[0],
			Size:         sz,
			ContentType:  rec[2],
			LastModified: rec[3],
		})
	}
	return objs, nil
}

func SaveObjects(dataDir, bucket string, objs []Object) error {
	f, err := os.Create(ObjectsFilePath(dataDir, bucket))
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, o := range objs {
		rec := []string{o.Key, strconv.FormatInt(o.Size, 10), o.ContentType, o.LastModified}
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func FindObject(objs []Object, key string) (int, Object) {
	for i, o := range objs {
		if o.Key == key {
			return i, o
		}
	}
	return -1, Object{}
}

func SaveFile(dataDir, bucket, key string, data []byte, contentType string) (Object, error) {
	if !isValidObjectKey(key) {
		return Object{}, ErrKeyInvalid
	}
	filePath := filepath.Join(dataDir, bucket, key)
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return Object{}, err
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return Object{}, err
	}
	objs, err := LoadObjects(dataDir, bucket)
	if err != nil {
		return Object{}, err
	}
	now := time.Now().Format(time.RFC3339)
	obj := Object{
		Key:          key,
		Size:         stat.Size(),
		ContentType:  contentType,
		LastModified: now,
	}
	idx, _ := FindObject(objs, key)
	if idx != -1 {
		objs[idx] = obj
	} else {
		objs = append(objs, obj)
	}
	if err := SaveObjects(dataDir, bucket, objs); err != nil {
		return Object{}, err
	}
	return obj, nil
}

func GetFile(dataDir, bucket, key string) ([]byte, Object, error) {
	objs, err := LoadObjects(dataDir, bucket)
	if err != nil {
		return nil, Object{}, err
	}
	idx, obj := FindObject(objs, key)
	if idx == -1 {
		return nil, Object{}, ErrNotFound
	}
	filePath := filepath.Join(dataDir, bucket, key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, Object{}, err
	}
	return data, obj, nil
}

func DeleteFile(dataDir, bucket, key string) error {
	objs, err := LoadObjects(dataDir, bucket)
	if err != nil {
		return err
	}
	idx, _ := FindObject(objs, key)
	if idx == -1 {
		return ErrNotFound
	}
	objs = append(objs[:idx], objs[idx+1:]...)
	if err := SaveObjects(dataDir, bucket, objs); err != nil {
		return err
	}
	filePath := filepath.Join(dataDir, bucket, key)
	return os.Remove(filePath)
}
