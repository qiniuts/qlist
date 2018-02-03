# qlist
list qiniu files in qiniu bucket or scan local file and proc the item

## edit config
edit your config 
```
{
  "access_key": "",
  "secret_key": "",
  "bucket": "",
  "done_records_fpath": "done_records.log",
  "proc_results_fpath": "proc_results.log",
  "worker_count": 10
}
```

## build

```
source env.sh
make env
make install
```

## usage

```
./qlist -h
Usage of ./qlist:
  -cfg_path string
    	 (default "cfg.json")
  -file_path string
    	 (default "keys.txt")
```
