package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AccessKey       string `json:"access_key"`
	SecretKey       string `json:"secret_key"`
	Bucket          string `json:"bucket"`
	FopQuery        string `json:"fop_query"`
	DoneRecordsPath string `json:"done_records_fpath"`
	ProcResultsPath string `json:"proc_results_fpath"`
	WorkerCount     int    `json:"worker_count"`
}

func LoadConfig(fpath string) (cfg Config, err error) {

	fh, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(fh)

	err = decoder.Decode(&cfg)
	return
}
