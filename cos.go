package tool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

type COSConf struct {
	SecretID   string `yaml:"secretid"`
	SecretKey  string `yaml:"secretkey"`
	Appid      string `yaml:"appid"`
	Region     string `yaml:"region"`
	Bucket     string `yaml:"bucket"`
	ExpireTime int64  `yaml:"expiretime"`
	TempDir    string `yaml:"tempdir"`
}
type CosDB int

func EnableCos(conf COSConf) error {
	var err error
	cosBase, err = newCosBase(conf.SecretID, conf.SecretKey, conf.Bucket, conf.Appid, conf.Region, conf.ExpireTime)
	if err != nil {
		return err
	}
	os.MkdirAll(conf.TempDir, 0755)
	os.Chmod(conf.TempDir, 0755)
	if err != nil {
		return err
	}
	return nil
}

var cosBase *CosBase

// CosBase 文件上传、下载操作与腾讯云 COS 服务交互需要的信息
type CosBase struct {
	SecretID             string
	SecretKey            string
	Bucket               string
	Appid                string
	Region               string
	CosURLExpireTimeNano int64
}

// NewCosBase ...
func newCosBase(secretID, secretKey, bucket, appid, region string, cosURLExpireTimeMs int64) (*CosBase, error) {
	return &CosBase{
		SecretID:             secretID,
		SecretKey:            secretKey,
		Bucket:               bucket,
		Appid:                appid,
		Region:               region,
		CosURLExpireTimeNano: cosURLExpireTimeMs * 1000 * 1000,
	}, nil
}

// GetCosFileURL 后台上传文件到COS服务上，返回预签名授权下载文件的链接和无签名的短链接durl
// filename 文件名字
// httpMethod为http.MethodGet时：生成文件下载URL，用于前端从COS服务器下载文件
// httpMethod为http.MethodPut时：生成文件上传URL，用于前端上传文件到COS服务器
func (cdb CosDB) GetCosFileURL(ctx context.Context, filename string, httpMethod string) (string, string) {
	c := cosBase.genClient(true)

	// Get presigned
	presignedURL, err := c.Object.GetPresignedURL(context.Background(), httpMethod, filename, cosBase.SecretID, cosBase.SecretKey, time.Duration(cosBase.CosURLExpireTimeNano), nil)
	if err != nil {
		return "", ""
	}

	durl := fmt.Sprintf("%s://%s%s", presignedURL.Scheme, presignedURL.Host, presignedURL.Path)
	return presignedURL.String(), durl
}

//ListDir 拉取bucket下的文件列表
func (cdb CosDB) ListDir(ctx context.Context) (*cos.BucketGetResult, error) {
	c := cosBase.genClient(false)

	opt := &cos.BucketGetOptions{
		MaxKeys: 1000,
	}
	v, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// UploadFileToCos 上传生成的excel文件到cos服务上
func (cdb CosDB) UploadFileToCos(ctx context.Context, filelocal, filename string) (string, string, error) {
	f, err := os.Open(filelocal)
	if err != nil {
		return "", "", err
	}
	s, err := f.Stat()
	if err != nil {
		return "", "", err
	}

	c := cosBase.genClient(false)
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength:      int(s.Size()),
			ContentDisposition: "attachment",
		},
	}
	_, err = c.Object.Put(context.Background(), filename, f, opt)
	if err != nil {
		return "", "", err
	}

	signUrl, dUrl := cdb.GetCosFileURL(ctx, filename, http.MethodGet)
	return signUrl, dUrl, nil
}

// UploadFileToCos 上传生成的excel文件到cos服务上
func (cdb CosDB) UploadFileToCosWithExpir(ctx context.Context, filelocal, filename string) (string, string, error) {
	f, err := os.Open(filelocal)
	if err != nil {
		return "", "", err
	}
	s, err := f.Stat()
	if err != nil {
		return "", "", err
	}

	c := cosBase.genClient(true)
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength:      int(s.Size()),
			ContentDisposition: "attachment",
		},
	}
	_, err = c.Object.Put(context.Background(), filename, f, opt)
	if err != nil {
		return "", "", err
	}

	signUrl, dUrl := cdb.GetCosFileURL(ctx, filename, http.MethodGet)
	return signUrl, dUrl, nil
}

// UploadFileToCos 上传生成的excel文件到cos服务上
func (cdb CosDB) UploadImageToCos(ctx context.Context, file io.Reader, filename string) (string, string, error) {
	c := cosBase.genClient(false)
	_, err := c.Object.Put(context.Background(), filename, file, nil)
	if err != nil {
		return "", "", err
	}

	signUrl, dUrl := cdb.GetCosFileURL(ctx, filename, http.MethodGet)
	return signUrl, dUrl, nil
}

// DownloadFileFromCos 从COS服务下载文件到本地
func (cdb CosDB) DownloadFileFromCos(ctx context.Context, filename, localpath string) error {
	c := cosBase.genClient(false)

	// 下载到本地的文件名和COS上的文件名一样，带有类似路径的
	_, err := c.Object.GetToFile(context.Background(), filename, localpath, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFileFromCos 删除cos上的文件
func (cdb CosDB) DeleteFileFromCos(ctx context.Context, filename string) error {
	c := cosBase.genClient(false)

	_, err := c.Object.Delete(context.Background(), filename)
	if err != nil {
		return err
	}

	return nil
}

func (a *CosBase) genClient(expired bool) *cos.Client {
	cosURL := "https://" + a.Bucket + "-" + a.Appid + ".cos." + a.Region + ".myqcloud.com"
	u, _ := url.Parse(cosURL)
	b := &cos.BaseURL{BucketURL: u}
	if expired {
		return cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  a.SecretID,
				SecretKey: a.SecretKey,
				Expire:    time.Duration(a.CosURLExpireTimeNano),
			},
		})
	}

	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  a.SecretID,
			SecretKey: a.SecretKey,
		},
	})
}
