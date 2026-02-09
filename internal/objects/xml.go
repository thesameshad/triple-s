package objects

import (
	"encoding/xml"
)

type ListBucketResult struct {
	XMLName  xml.Name    `xml:"ListBucketResult"`
	Contents []ObjectXML `xml:"Contents"`
}

type ObjectXML struct {
	Key          string `xml:"Key"`
	Size         int64  `xml:"Size"`
	ContentType  string `xml:"ContentType"`
	LastModified string `xml:"LastModified"`
}

func ObjectsToXML(objs []Object) ([]byte, error) {
	res := ListBucketResult{}
	for _, o := range objs {
		res.Contents = append(res.Contents, ObjectXML{
			Key: o.Key, Size: o.Size, ContentType: o.ContentType, LastModified: o.LastModified,
		})
	}
	return xml.MarshalIndent(res, "", "  ")
}
