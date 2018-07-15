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
	Prefix    string `json:"prefix"`

	//dest proc
	FopQuery        string `json:"fop_query"`
	ProcResultsPath string `json:"proc_results_fpath"`

	//concurency num
	WorkerCount int `json:"worker_count"`

	ReqHeaderHost string `json:"req_header_host"`

	ExtraParams map[string]map[string]interface{} `json:"extra_params"`
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

func (c Config) GetExtraParam(key1, key2 string) (param interface{}, ok bool) {

	params, ok := c.ExtraParams[key1]
	if !ok {
		return
	}

	param, ok = params[key2]
	return
}
