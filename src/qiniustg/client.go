package qiniustg

import (
	"config"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type Client struct {
	config.Config
}

func NewClient(cfg config.Config) *Client {
	return &Client{cfg}
}

func (c *Client) BucketMgr() *storage.BucketManager {

	mac := qbox.NewMac(c.AccessKey, c.SecretKey)
	stgCfg := storage.Config{
		UseHTTPS: false,
	}

	return storage.NewBucketManager(mac, &stgCfg)
}
