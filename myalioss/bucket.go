package myalioss

import (
	"bytes"
	"os"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Bucket 存储桶
type Bucket struct {
	Bucket *oss.Bucket
	Info   *BucketInfo
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
func CreateOSSBucket(endPoint, accessKeyID, accessKeySecret, bucketName string) (*Bucket, error) {
	// 客户端连接
	client, err := oss.New(endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	// 判断存储空间是否
	isExist, err := client.IsBucketExist(bucketName)
	if err != nil {
		return nil, err
	}

	// 创建一个存储桶
	if !isExist {
		// 创建存储空间，并设置存储类型为标准访问oss.StorageStandard、读写权限ACL为公共读oss.ACLPublicRead、数据容灾类型为本地冗余存储oss.RedundancyLRS
		err = client.CreateBucket(bucketName, oss.StorageClass(oss.StorageStandard), oss.ACL(oss.ACLPublicRead), oss.RedundancyType(oss.RedundancyLRS)) // 名字只能包含小写字母、数字、`-`，桶名称易重复
		if err != nil {
			return nil, err
		}
	}

	// 获取存储空间的信息
	res, err := client.GetBucketInfo(bucketName)
	if err != nil {
		return nil, err
	}

	// fmt.Println("BucketInfo.Location: ", res.BucketInfo.Location)
	// fmt.Println("BucketInfo.CreationDate: ", res.BucketInfo.CreationDate)
	// fmt.Println("BucketInfo.ACL: ", res.BucketInfo.ACL)
	// fmt.Println("BucketInfo.Owner: ", res.BucketInfo.Owner)
	// fmt.Println("BucketInfo.StorageClass: ", res.BucketInfo.StorageClass)
	// fmt.Println("BucketInfo.RedundancyType: ", res.BucketInfo.RedundancyType)
	// fmt.Println("BucketInfo.ExtranetEndpoint: ", res.BucketInfo.ExtranetEndpoint)
	// fmt.Println("BucketInfo.IntranetEndpoint: ", res.BucketInfo.IntranetEndpoint)

	_bucket := &Bucket{
		Info: &BucketInfo{
			Name:                   res.BucketInfo.Name,
			Location:               res.BucketInfo.Location,
			CreationDate:           res.BucketInfo.CreationDate,
			ExtranetEndpoint:       res.BucketInfo.ExtranetEndpoint,
			IntranetEndpoint:       res.BucketInfo.IntranetEndpoint,
			ACL:                    res.BucketInfo.ACL,
			RedundancyType:         res.BucketInfo.RedundancyType,
			StorageClass:           res.BucketInfo.StorageClass,
			Versioning:             res.BucketInfo.Versioning,
			TransferAcceleration:   res.BucketInfo.TransferAcceleration,
			CrossRegionReplication: res.BucketInfo.CrossRegionReplication,
		},
	}

	// 获取存储空间。
	_bucket.Bucket, err = client.Bucket(bucketName)
	if err != nil {
		return _bucket, err
	}

	return _bucket, err
}

// GetObjectToFile 下载到本地文件
func (b *Bucket) GetObjectToFile(objectKey string, filePath string) error {
	return b.Bucket.GetObjectToFile(objectKey, filePath)
}

// PutObjectFromFileName 上传本地文件
func (b *Bucket) PutObjectFromFileName(from string, to string) error {
	return b.Bucket.PutObjectFromFile(to, from)
}

// PutObjectFromFile 上传文件流
func (b *Bucket) PutObjectFromFile(from string, to string) error {
	file, err := os.Open(from)
	if err != nil {
		return err
	}
	defer file.Close()

	return b.Bucket.PutObject(to, file)
}

// PutObjectFromReader 上传文件流
func (b *Bucket) PutObjectFromBytes(from []byte, to string) error {
	reader := bytes.NewReader(from)
	return b.Bucket.PutObject(to, reader)
}
