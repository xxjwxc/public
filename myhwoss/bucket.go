package myhwoss

import (
	"bytes"
	"fmt"
	"os"
	"time"

	obs "github.com/xxjwxc/public/myhwoss/obs"
)

// Bucket 存储桶
type Bucket struct {
	ObsClient  *obs.ObsClient
	BucketName string
}

// BucketInfo 存储桶的信息
type BucketInfo struct {
	Name                   string    `xml:"Name"`                    // 桶名
	Location               string    `xml:"Location"`                // Bucket所在的地域
	CreationDate           time.Time `xml:"CreationDate"`            // Bucket的创建时间，格式为UTC时间
	ExtranetEndpoint       string    `xml:"ExtranetEndpoint"`        // Bucket的外网域名
	IntranetEndpoint       string    `xml:"IntranetEndpoint"`        // 同地域ECS访问Bucket的内网域名
	ACL                    string    `xml:"AccessControlList>Grant"` // Bucket读写权限（ACL）信息的容器
	RedundancyType         string    `xml:"DataRedundancyType"`      // Bucket的数据容灾类型
	StorageClass           string    `xml:"StorageClass"`            // Bucket的存储类型
	Versioning             string    `xml:"Versioning"`              // Bucket的版本控制状态。有效值：Enabled、Suspended
	TransferAcceleration   string    `xml:"TransferAcceleration"`    // 显示Bucket的传输加速状态。有效值：Enabled、Disabled
	CrossRegionReplication string    `xml:"CrossRegionReplication"`  // 显示Bucket的跨区域复制状态。有效值：Enabled、Disabled
}

// CreateOSSBucket 获取OSS对象存储桶
func CreateOSSBucket(endPoint, ak, sk, bucketName string) (*Bucket, error) {
	// 客户端连接
	obsClient, err := obs.New(ak, sk, endPoint)
	if err != nil {
		panic(err)
	}

	// 判断存储空间是否
	output, err := obsClient.HeadBucket(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if _, ok := err.(obs.ObsError); ok {
			// fmt.Println(obsError.StatusCode)
			// 创建存储空间，并设置存储类型为标准访问oss.StorageStandard、读写权限ACL为公共读oss.ACLPublicRead、数据容灾类型为本地冗余存储oss.RedundancyLRS
			if err != nil {
				return nil, err
			}
			input := &obs.CreateBucketInput{}
			input.Bucket = bucketName
			// input.StorageClass = obs.StorageClassWarm
			input.ACL = obs.AclPublicRead
			_, err := obsClient.CreateBucket(input)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	return &Bucket{
		ObsClient:  obsClient,
		BucketName: bucketName,
	}, err
}

// GetObjectToFile 下载到本地文件
func (b *Bucket) GetObjectToFile(objectKey string, filePath string) error {
	input := &obs.GetObjectInput{}
	input.Bucket = b.BucketName
	input.Key = objectKey
	output, err := b.ObsClient.GetObject(input, obs.WithProgress(&ObsProgressListener{}))
	if err != nil {
		return err
	}
	defer output.Body.Close()

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	f, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 10240)
	var readErr error
	var readCount int
	for {
		readCount, readErr = output.Body.Read(p)
		if readCount > 0 {
			f.Write(p[:readCount])
		}
		if readErr != nil {
			break
		}
	}

	return nil
}

// PutObjectFromFileName 上传本地文件
func (b *Bucket) PutObjectFromFileName(from string, to string) error {
	input := &obs.PutFileInput{}
	input.Bucket = b.BucketName
	input.Key = to
	input.SourceFile = from
	_, err := b.ObsClient.PutFile(input)
	return err
}

// PutObjectFromFile 上传文件流
func (b *Bucket) PutObjectFromFile(from string, to string) error {
	return b.PutObjectFromFileName(from, to)
}

// PutObjectFromReader 上传文件流
func (b *Bucket) PutObjectFromBytes(from []byte, to string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = b.BucketName
	input.Key = to
	input.Body = bytes.NewReader(from)
	_, err := b.ObsClient.PutObject(input, obs.WithProgress(&ObsProgressListener{}))
	return err
}

// IsObjectExist 判断文件是否存在
func (b *Bucket) IsObjectExist(filename string) (bool, error) {
	input := &obs.GetObjectMetadataInput{}
	input.Bucket = b.BucketName
	input.Key = filename
	_, err := b.ObsClient.GetObjectMetadata(input)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// DeleteObject 删除文件
func (b *Bucket) DeleteObject(filenames []string) ([]string, error) {
	input := &obs.DeleteObjectsInput{}
	input.Bucket = b.BucketName
	for _, v := range filenames {
		input.Objects = append(input.Objects, obs.ObjectToDelete{Key: v})
	}
	output, err := b.ObsClient.DeleteObjects(input)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, deleted := range output.Deleteds {
		out = append(out, deleted.Key)
	}

	return out, nil
}

// // ListObjects 列举文件
// func (b *Bucket) ListObjects(filenames []string) ([]string, error) {
// 	lsRes, err := b.Bucket.ListObjects(oss.Marker(marker))
// 	return res.DeletedObjects, err
// }

// SetObjectACL 设置文件的访问权限
func (b *Bucket) SetObjectACL(objectKey string, _type []obs.Grant) error {
	input := &obs.SetObjectAclInput{}
	input.Bucket = b.BucketName
	input.Key = objectKey
	input.ACL = obs.AclPublicRead
	input.Grants = _type
	_, err := b.ObsClient.SetObjectAcl(input)
	return err
}

// GetObjectACL 获取文件的访问权限
func (b *Bucket) GetObjectACL(objectKey string) ([]obs.Grant, error) {
	input := &obs.GetObjectAclInput{}
	input.Bucket = b.BucketName
	input.Key = objectKey
	output, err := b.ObsClient.GetObjectAcl(input)
	if err != nil {
		return nil, err
	}

	return output.Grants, nil
}

// GetObjectMeta 获取文件元信息
func (b *Bucket) GetObjectMeta(objectName string) (map[string][]string, error) {
	input := &obs.GetObjectMetadataInput{}
	input.Bucket = b.BucketName
	input.Key = objectName
	output, err := b.ObsClient.GetObjectMetadata(input)
	if err != nil {
		return nil, err
	}

	mp := make(map[string][]string)
	mp["ETag"] = append(mp["ETag"], output.ETag)
	mp["ContentType"] = append(mp["ContentType"], output.ContentType)

	return mp, nil
}
