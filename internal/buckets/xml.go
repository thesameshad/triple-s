package buckets

import (
	"encoding/xml"
)

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets BucketsElem `xml:"Buckets"`
}

type BucketsElem struct {
	Bucket []BucketElem `xml:"Bucket"`
}

type BucketElem struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
	LastModified string `xml:"LastModified"`
}

func BucketsToXML(buckets []Bucket) ([]byte, error) {
	lst := ListAllMyBucketsResult{
		Buckets: BucketsElem{},
	}
	for _, b := range buckets {
		if b.Status != "active" {
			continue
		}
		lst.Buckets.Bucket = append(lst.Buckets.Bucket, BucketElem{
			Name: b.Name,
			CreationDate: b.CreationTime,
			LastModified: b.LastModified,
		})
	}
	return xml.MarshalIndent(lst, "", "  ")
}
