package buckets

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Bucket struct {
	Name         string
	CreationTime string
	LastModified string
	Status       string
}

var (
	bucketNameRegexp     = regexp.MustCompile(`^[a-z0-9][a-z0-9-.]{1,61}[a-z0-9]$`)
	ErrBucketExists      = errors.New("Bucket already exists")
	ErrBucketNotFound    = errors.New("Bucket not found")
	ErrInvalidBucketName = errors.New("Invalid bucket name")
	ErrBucketNotEmpty    = errors.New("Bucket is not empty")
)

func isValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}
	if !bucketNameRegexp.MatchString(name) {
		return false
	}
	return true
}

func BucketsFilePath(dataDir string) string {
	return filepath.Join(dataDir, "buckets.csv")
}

func LoadBuckets(dataDir string) ([]Bucket, error) {
	f, err := os.Open(BucketsFilePath(dataDir))
	if os.IsNotExist(err) {
		return []Bucket{}, nil
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
	var buckets []Bucket
	for _, rec := range records {
		if len(rec) < 4 {
			continue
		}
		buckets = append(buckets, Bucket{
			Name:         rec[0],
			CreationTime: rec[1],
			LastModified: rec[2],
			Status:       rec[3],
		})
	}
	return buckets, nil
}

func SaveBuckets(dataDir string, buckets []Bucket) error {
	f, err := os.Create(BucketsFilePath(dataDir))
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, b := range buckets {
		rec := []string{b.Name, b.CreationTime, b.LastModified, b.Status}
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func FindBucket(buckets []Bucket, name string) (int, Bucket) {
	for i, b := range buckets {
		if b.Name == name {
			return i, b
		}
	}
	return -1, Bucket{}
}

func CreateBucket(dataDir, name string) (Bucket, error) {
	if !isValidBucketName(name) {
		return Bucket{}, ErrInvalidBucketName
	}
	buckets, err := LoadBuckets(dataDir)
	if err != nil {
		return Bucket{}, err
	}
	if idx, _ := FindBucket(buckets, name); idx != -1 {
		return Bucket{}, ErrBucketExists
	}
	now := time.Now().Format(time.RFC3339)
	bucket := Bucket{
		Name:         name,
		CreationTime: now,
		LastModified: now,
		Status:       "active",
	}
	buckets = append(buckets, bucket)
	if err := SaveBuckets(dataDir, buckets); err != nil {
		return Bucket{}, err
	}
	if err := os.MkdirAll(filepath.Join(dataDir, name), 0o755); err != nil {
		return Bucket{}, err
	}
	return bucket, nil
}

func DeleteBucket(dataDir, name string) error {
	buckets, err := LoadBuckets(dataDir)
	if err != nil {
		return err
	}
	idx, _ := FindBucket(buckets, name)
	if idx == -1 {
		return ErrBucketNotFound
	}
	objDir := filepath.Join(dataDir, name)
	files, err := os.ReadDir(objDir)
	if err == nil && len(files) > 0 {
		return ErrBucketNotEmpty
	}
	buckets = append(buckets[:idx], buckets[idx+1:]...)
	if err := SaveBuckets(dataDir, buckets); err != nil {
		return err
	}
	os.RemoveAll(objDir)
	return nil
}
