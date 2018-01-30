package qiniustg

import "config"

type Client struct {
	config.Config
}

func NewClient(cfg config.Config) *Client {
	return &Client{cfg}
}