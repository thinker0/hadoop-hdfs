package main

import (
	"encoding/json"
	"fmt"
	"github.com/thinker0/hadoop-hdfs/v2/pkg/hadoop/hdfs/fsimage"
	"go.uber.org/zap"
	"log"
	"os"
)

// Global Variables
var pathCounter int = 0

func logIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	sampleJSON := []byte(`{
       "level" : "info",
       "encoding": "json",
       "outputPaths":["stdout", "log.log"],
       "errorOutputPaths":["stderr"],
       "encoderConfig": {
           "messageKey":"message",
           "levelKey":"level",
           "levelEncoder":"lowercase"
       }
   }`)
	var cfg zap.Config

	if err := json.Unmarshal(sampleJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	fileName := os.Args[1]
	fInfo, err := os.Stat(fileName)
	logIfErr(err)

	f, err := os.Open(fileName)
	logIfErr(err)

	fileLength := fInfo.Size()
	fSummaryLength := fsimage.DecodeFileSummaryLength(fileLength, f)
	sectionMap := fsimage.ParseFileSummary(f, fileLength, fSummaryLength)

	inodeSectionInfo := sectionMap["INODE"]
	inodeNames, entityCount := fsimage.ParseInodeSection(inodeSectionInfo, f)

	inodeDirectorySectionInfo := sectionMap["INODE_DIR"]
	parChildrenMap := fsimage.ParseInodeDirectorySection(inodeDirectorySectionInfo, f)
	fmt.Println("Root INode ID: ", fsimage.ROOT_INODE_ID)
	fmt.Println("Total Number of Files: ", entityCount.Files)
	fmt.Println("Total Number of Directories: ", entityCount.Directories)
	fmt.Println("Total Number of Symlinks: ", entityCount.Symlinks)
	rootTreeNode := fsimage.FindChildren(parChildrenMap, inodeNames, fsimage.ROOT_INODE_ID)
	cnt := fsimage.CountTreeNodes(*rootTreeNode)
	paths := make([]string, cnt)
	var index int = 0
	for _, child := range rootTreeNode.Children {
		fmt.Println("Processing " + string(child.Name))
		fsimage.ConstructPath("", child, &paths, &index)
	}
	p := 0
	for _, c := range paths {
		if len(c) != 0 {
			p++
			log.Println(c)
		}
	}
	fmt.Println("No of Paths: ", p)
	fmt.Println("First 10 paths")
	fmt.Println(paths[:10])

	fmt.Println("Parse further")
	fmt.Println(sectionMap)
	f.Close()
}
