# qlist
list qiniu files in qiniu bucket or scan local file and proc the item

## edit config
edit your config 
```
{
  "src": "localstg",
  "to_do_records_path": "to_do_records_path.txt",
  "access_key": "",
  "secret_key": "",
  "bucket": "",
  "fop_query": "?imageView2/2/w/720|qpolitician",
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

```

## batch change file type
```
./qlist -cfg_path cfg.json chtype
```

## batch change file status
```
./qlist -cfg_path cfg.json chstatus
```

## proc file by fop
```
./qlist -cfg_path cfg.json req
```

## list bucket files
```
./qlist -cfg_path cfg.json bucketlist
```

## async fetch url
url file format: 
```
<url> \t <key> \t <md5>
```

run:
```
./qlist -cfg_path cfg.json async_fetch
```
