package config

import (
	"os"
	"encoding/json"
)

type Config struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket string `json:"bucket"`
	DoneRecordPath string `json:"success_record_path"`
	WorkerCount int `json:"worker_count"`
}

func LoadConfig(fpath string) (cfg Config, err error)  {

	fh, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(fh)

	err = decoder.Decode(&cfg)
	return
}
