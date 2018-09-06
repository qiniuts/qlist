package qiniustg

import (
	"config"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type QNClient struct {
	config.Config
}

func NewQNClient(cfg config.Config) *QNClient {
	return &QNClient{cfg}
}

func (c *QNClient) BucketMgr() *storage.BucketManager {

	mac := qbox.NewMac(c.AccessKey, c.SecretKey)
	stgCfg := storage.Config{
		UseHTTPS: false,
	}

	return storage.NewBucketManager(mac, &stgCfg)
}
