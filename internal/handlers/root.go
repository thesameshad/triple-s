package handlers

import (
	"io"
	"net/http"
	"strings"
	"triple-s/internal/buckets"
	"triple-s/internal/objects"
)

var dataDir string

func SetDataDir(dir string) {
	dataDir = dir
}

func Root(w http.ResponseWriter, r *http.Request) {
	apiURL := r.URL.Path
	parts := strings.Split(strings.Trim(apiURL, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		if r.Method == http.MethodGet {
			listBuckets(w)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(parts) == 1 {
		bucketName := parts[0]
		if r.Method == http.MethodPut {
			createBucket(w, bucketName)
			return
		} else if r.Method == http.MethodDelete {
			deleteBucket(w, bucketName)
			return
		} else if r.Method == http.MethodGet {
			listObjectsXML(w, bucketName)
			return
		}
	}
	if len(parts) == 2 {
		bucketName := parts[0]
		objectKey := parts[1]
		if r.Method == http.MethodPut {
			uploadObject(w, r, bucketName, objectKey)
			return
		} else if r.Method == http.MethodGet {
			getObject(w, bucketName, objectKey)
			return
		} else if r.Method == http.MethodDelete {
			deleteObject(w, bucketName, objectKey)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func listBuckets(w http.ResponseWriter) {
	bks, err := buckets.LoadBuckets(dataDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	xmlBytes, _ := buckets.BucketsToXML(bks)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(xmlBytes)
}

func listObjectsXML(w http.ResponseWriter, bucket string) {
	bks, err := buckets.LoadBuckets(dataDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	_, bk := buckets.FindBucket(bks, bucket)
	if bk.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<Error>NoSuchBucket</Error>"))
		return
	}
	objs, err := objects.LoadObjects(dataDir, bucket)
	if err != nil {
		objs = []objects.Object{}
	}
	xmlBytes, _ := objects.ObjectsToXML(objs)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(xmlBytes)
}

func createBucket(w http.ResponseWriter, name string) {
	bucket, err := buckets.CreateBucket(dataDir, name)
	if err != nil {
		if err == buckets.ErrInvalidBucketName {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<Error>InvalidBucketName</Error>"))
			return
		}
		if err == buckets.ErrBucketExists {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("<Error>BucketAlreadyExists</Error>"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<Bucket>" + bucket.Name + "</Bucket>"))
}

func deleteBucket(w http.ResponseWriter, name string) {
	err := buckets.DeleteBucket(dataDir, name)
	if err != nil {
		if err == buckets.ErrBucketNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("<Error>NoSuchBucket</Error>"))
			return
		}
		if err == buckets.ErrBucketNotEmpty {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("<Error>BucketNotEmpty</Error>"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func uploadObject(w http.ResponseWriter, r *http.Request, bucket, key string) {
	bks, err := buckets.LoadBuckets(dataDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	_, bk := buckets.FindBucket(bks, bucket)
	if bk.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<Error>NoSuchBucket</Error>"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<Error>InvalidRequest</Error>"))
		return
	}
	ctype := r.Header.Get("Content-Type")
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	obj, err := objects.SaveFile(dataDir, bucket, key, body, ctype)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<Error>InvalidObjectKey</Error>"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<ObjectKey>" + obj.Key + "</ObjectKey>"))
}

func getObject(w http.ResponseWriter, bucket, key string) {
	bks, err := buckets.LoadBuckets(dataDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	_, bk := buckets.FindBucket(bks, bucket)
	if bk.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<Error>NoSuchBucket</Error>"))
		return
	}
	file, obj, err := objects.GetFile(dataDir, bucket, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<Error>NoSuchKey</Error>"))
		return
	}
	w.Header().Set("Content-Type", obj.ContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func deleteObject(w http.ResponseWriter, bucket, key string) {
	bks, err := buckets.LoadBuckets(dataDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<Error>InternalError</Error>"))
		return
	}
	_, bk := buckets.FindBucket(bks, bucket)
	if bk.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<Error>NoSuchBucket</Error>"))
		return
	}
	err = objects.DeleteFile(dataDir, bucket, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<Error>NoSuchKey</Error>"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
