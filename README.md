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
`GOOS=linux GOARCH=amd64 go build cmd/hdfs-fsimage-parse/*.go`

## Run

`go run cmd/hdfs-fsimage-parse/*.go fsimage_0001`


## Output
```
main.go:180: Section Name: NS_INFO:name:"NS_INFO"  length:24  offset:8
main.go:180: Section Name: INODE:name:"INODE"  length:1925  offset:32
main.go:180: Section Name: INODE_DIR:name:"INODE_DIR"  length:157  offset:1957
main.go:180: Section Name: FILES_UNDERCONSTRUCTION:name:"FILES_UNDERCONSTRUCTION"  length:0  offset:2114
main.go:180: Section Name: SNAPSHOT:name:"SNAPSHOT"  length:5  offset:2114
main.go:180: Section Name: SNAPSHOT_DIFF:name:"SNAPSHOT_DIFF"  length:9  offset:2119
main.go:180: Section Name: INODE_REFERENCE:name:"INODE_REFERENCE"  length:0  offset:2128
main.go:180: Section Name: SECRET_MANAGER:name:"SECRET_MANAGER"  length:9  offset:2128
main.go:180: Section Name: CACHE_MANAGER:name:"CACHE_MANAGER"  length:7  offset:2137
main.go:180: Section Name: STRING_TABLE:name:"STRING_TABLE"  length:53  offset:2144
Root INode ID:  16385
Total Number of Files:  16
Total Number of Directories:  14
Total Number of Symlinks:  0
Processing datalake
Processing test1
Processing test2
Processing test3
Processing test_2KiB.img
Processing user
main.go:90: /datalake/asset1
main.go:90: /datalake/asset2/test_1KiB.img
main.go:90: /datalake/asset2/test_2MiB.img
main.go:90: /datalake/asset3/subasset1/test_2MiB.img
main.go:90: /datalake/asset3/subasset2/test_2MiB.img
main.go:90: /datalake/asset3/test_2MiB.img
main.go:90: /test1
main.go:90: /test2
main.go:90: /test3/foo/bar/test_20MiB.img
main.go:90: /test3/foo/bar/test_2MiB.img
main.go:90: /test3/foo/bar/test_40MiB.img
main.go:90: /test3/foo/bar/test_4MiB.img
main.go:90: /test3/foo/bar/test_5MiB.img
main.go:90: /test3/foo/bar/test_80MiB.img
main.go:90: /test3/foo/test_1KiB.img
main.go:90: /test3/foo/test_20MiB.img
main.go:90: /test3/test.img
main.go:90: /test3/test_160MiB.img
main.go:90: /test_2KiB.img
main.go:90: /user/mm
No of Paths:  20
First 10 paths
[/datalake/asset1 /datalake/asset2/test_1KiB.img /datalake/asset2/test_2MiB.img /datalake/asset3/subasset1/test_2MiB.img /datalake/asset3/subasset2/test_2MiB.img /datalake/asset3/test_2MiB.img /test1 /test2 /test3/foo/bar/test_20MiB.img /test3/foo/bar/test_2MiB.img]
Parse further
map[CACHE_MANAGER:name:"CACHE_MANAGER"  length:7  offset:2137 FILES_UNDERCONSTRUCTION:name:"FILES_UNDERCONSTRUCTION"  length:0  offset:2114 INODE:name:"INODE"  length:1925  offset:32 INODE_DIR:name:"INODE_DIR"  length:157  offset:1957 INODE_REFERENCE:name:"INODE_REFERENCE"  length:0  offset:2128 NS_INFO:name:"NS_INFO"  length:24  offset:8 SECRET_MANAGER:name:"SECRET_MANAGER"  length:9  offset:2128 SNAPSHOT:name:"SNAPSHOT"  length:5  offset:2114 SNAPSHOT_DIFF:name:"SNAPSHOT_DIFF"  length:9  offset:2119 STRING_TABLE:name:"STRING_TABLE"  length:53  offset:2144]
```

