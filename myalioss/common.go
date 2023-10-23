package myalioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// IsObjectExist 判断文件是否存在
func (b *Bucket) IsObjectExist(filename string) (bool, error) {
	return b.Bucket.IsObjectExist(filename)
}

// DeleteObject 删除文件
func (b *Bucket) DeleteObject(filenames []string) ([]string, error) {
	res, err := b.Bucket.DeleteObjects(filenames)
	return res.DeletedObjects, err
}

// // ListObjects 列举文件
// func (b *Bucket) ListObjects(filenames []string) ([]string, error) {
// 	lsRes, err := b.Bucket.ListObjects(oss.Marker(marker))
// 	return res.DeletedObjects, err
// }

// SetObjectACL 设置文件的访问权限
func (b *Bucket) SetObjectACL(_type oss.ACLType) error {
	return b.Bucket.SetObjectACL(b.Info.Name, _type)
}

// GetObjectACL 获取文件的访问权限
func (b *Bucket) GetObjectACL() (string, error) {
	aclRes, err := b.Bucket.GetObjectACL(b.Info.Name)
	return aclRes.ACL, err
}

// GetObjectMeta 获取文件元信息
func (b *Bucket) GetObjectMeta(objectName string) (map[string][]string, error) {
	props, err := b.Bucket.GetObjectMeta(objectName)
	if err != nil {
		return nil, err
	}

	return props, nil
}
