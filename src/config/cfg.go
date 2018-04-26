package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Src string `json:"src"` //qiniustg|localstg

	//localstg
	ToDoRecordsPath string `json:"to_do_records_path"`

	//qiniustg
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`

	//dest proc
	FopQuery        string `json:"fop_query"`
	ProcResultsPath string `json:"proc_results_fpath"`

	//concurency num
	WorkerCount int `json:"worker_count"`

	ReqHeaderHost string `json:"req_header_host"`
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

func (c Config) IsQiniuSrc() bool {
	return c.Src == "qiniustg"
}
