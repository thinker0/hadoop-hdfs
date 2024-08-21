# hdfs-fsimage-parse-go
A Golang port of https://github.com/sebinjohn/hdfs-fsimage-parse

## First Steps

```
$ protoc \
	--go_out=pkg/hadoop/common \
	--go_opt=paths=source_relative \
	-Ipkg/hadoop/common \
	pkg/hadoop/common/Security.proto

$ protoc \
    --proto_path=pkg/hadoop/hdfs \
    --proto_path=pkg/hadoop/common \
	--go_out=./pkg/hadoop/hdfs \
	-Ipkg/hadoop/hdfs -Ipkg/hadoop/common \
	--go_opt=paths=source_relative \
	 pkg/hadoop/hdfs/*.proto

$ protoc \
	--go_out=./pkg/hadoop/hdfs \
	-Ipkg/hadoop/hdfs \
	-Ipkg/hadoop/common \
	--go_opt=paths=source_relative \
	pkg/hadoop/hdfs/fsimage/fsimage.proto
```

## Build
for a linux x64 env
`GOOS=linux GOARCH=amd64 go build *.go`

## Run

`go run *.go <path to hdfs fsimage>`


## Output
```

Root INode ID:  16385
Total Number of Files:  40430334
Total Number of Directories:  11691564
Total Number of Symlinks:  0
Processing archive
Processing benchmarks
Processing dev
Processing jobtracker
Processing system
Processing test
Processing tmp
Processing user
No of Paths:  40911643
First 10 paths
[/user/xxx/path1 /user/some/path2 /tmp/11/path3 /tmp/path4 /benchmark/path5 /dev/path6 /jobtracker/path7 path8 /test/1111/path9 /dev]
Parse further
map[INODE:name:"INODE" length:4228334479 offset:46  INODE_DIR:name:"INODE_DIR" length:342395119 offset:4228334525  SECRET_MANAGER:name:"SECRET_MANAGER" length:4342600 offset:4570782021  CACHE_MANAGER:name:"CACHE_MANAGER" length:7 offset:4575124621  NS _INFO:name:"NS_INFO" length:38 offset:8  SNAPSHOT:name:"SNAPSHOT" length:5 offset:4570748819  SNAPSHOT_DIFF:name:"SNAPSHOT_DIFF" length:33197 offset:4570748824  INODE_REFERENCE:name:"INODE_REFERENCE" length:0 offset:4570782021  STRING_TABLE:name:"STRI
NG_TABLE" length:570 offset:4575124628  FILES_UNDERCONSTRUCTION:name:"FILES_UNDERCONSTRUCTION" length:19175 offset:4570729644 ]
```

