package localstg

import "config"

type QNClient struct {
	config.Config
}

func NewQNClient(cfg config.Config) *QNClient {
	return &QNClient{cfg}
}
