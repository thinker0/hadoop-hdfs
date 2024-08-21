package fsimage

import (
	"bytes"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

const (
	FILE      = iota
	DIRECTORY = iota
	SYMLINK   = iota

	FILE_SUM_BYTES = 4
	ROOT_INODE_ID  = 16385
)

type InodeId uint64
type ChildId uint64
type HDFSFileName string
type NameCount uint32

type ChildrenCount int
type ParentId uint64

type INodeTree struct {
	INode
	Children []*INodeTree
}

type EntityCount struct {
	Files       uint32
	Directories uint32
	Symlinks    uint32
}

type INode struct {
	Name []byte
	Id   InodeId
	Type int
}

func (i ChildrenCount) String() string {
	return strconv.Itoa(int(i))
}

func (i ParentId) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

type NameCountPair struct {
	Name  HDFSFileName
	Count NameCount
}

func (p NameCountPair) String() string {
	return string(p.Name) + " " + strconv.FormatUint(uint64(p.Count), 10)
}

type NameCountPairList []NameCountPair

func (p NameCountPairList) Len() int { return len(p) }

func (p NameCountPairList) Less(i, j int) bool { return p[i].Count < p[j].Count }

func (p NameCountPairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SortByNameCount(m map[HDFSFileName]NameCount) NameCountPairList {
	pl := make(NameCountPairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = NameCountPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type ChildrenCountPair struct {
	Parid      ParentId
	ChildCount ChildrenCount
}

type ChildrenCountPairList []ChildrenCountPair

func (p ChildrenCountPair) String() string {
	return p.Parid.String() + " " + p.ChildCount.String()
}

func (p ChildrenCountPairList) Len() int { return len(p) }

func (p ChildrenCountPairList) Less(i, j int) bool { return p[i].ChildCount < p[j].ChildCount }

func (p ChildrenCountPairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SortByChildCount(data map[ParentId]ChildrenCount) ChildrenCountPairList {
	pl := make(ChildrenCountPairList, len(data))
	i := 0
	for k, v := range data {
		pl[i] = ChildrenCountPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func MinINodeID(inodeIds []InodeId) InodeId {
	min := inodeIds[0]
	for _, v := range inodeIds {
		if v < min {
			min = v
		}
	}
	return min
}

func CountTreeNodes(rootNode INodeTree) uint32 {
	var i uint32 = 1
	for _, child := range rootNode.Children {
		i++
		countInSubTree(*child, &i)
	}
	log.Println("Total Number of Nodes: ", i)
	return i
}

func countInSubTree(node INodeTree, counter *uint32) {
	children := node.Children
	if children == nil || len(children) == 0 {
		(*counter)++
		return
	}
	for _, c := range node.Children {
		(*counter)++
		countInSubTree(*c, counter)
	}
}

func ConstructPath(constructedPath string, node *INodeTree, paths *[]string, index *int) {
	children := node.Children
	if children == nil || len(children) == 0 {
		log.Println("Processing " + constructedPath + "/" + string(node.Name))
		(*paths)[*index] = constructedPath + "/" + string(node.Name)
		(*index)++
	} else {
		for _, child := range children {
			log.Println("Processing " + constructedPath + "/" + string(node.Name))
			ConstructPath(constructedPath+"/"+string(node.Name), child, paths, index)
		}
	}
}

func FindChildren(parChildrenMap map[ParentId][]uint64, inodeNames map[InodeId]INode, curInodeId InodeId) *INodeTree {
	var children []uint64
	children, ok1 := parChildrenMap[ParentId(curInodeId)]
	inode, _ := inodeNames[curInodeId]
	numberOfChildren := len(children)

	if !ok1 || numberOfChildren == 0 {
		return &INodeTree{inode, nil}
	}

	refs := make([]*INodeTree, numberOfChildren)
	for i, child := range children {
		refs[i] = FindChildren(parChildrenMap, inodeNames, InodeId(child))
	}
	return &INodeTree{inode, refs}
}

func logIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DecodeFileSummaryLength(fileLength int64, imageFile *os.File) int32 {
	var (
		fSumLenBytes   = make([]byte, FILE_SUM_BYTES)
		fSummaryLength int32
	)
	fileSummaryLengthStart := fileLength - FILE_SUM_BYTES
	bReader := bytes.NewReader(fSumLenBytes)
	_, err := imageFile.ReadAt(fSumLenBytes, fileSummaryLengthStart)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
	if err = binary.Read(bReader, binary.BigEndian, &fSummaryLength); err != nil {
		log.Fatal(err)
	}
	return fSummaryLength
}

func ParseFileSummary(imageFile *os.File, fileLength int64, fSummaryLength int32) map[string]*FileSummary_Section {
	var (
		sectionMap               = make(map[string]*FileSummary_Section)
		fileSummary *FileSummary = &FileSummary{}
	)
	// last 4 bytes says how many bytes should be read from end to get the FileSummary message
	fSummaryLength64 := int64(fSummaryLength)
	readAt := fileLength - fSummaryLength64 - FILE_SUM_BYTES

	fSummaryBytes := make([]byte, fSummaryLength)
	_, err := imageFile.ReadAt(fSummaryBytes, readAt)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

	_, c := binary.Uvarint(fSummaryBytes)
	if c <= 0 {
		log.Fatal("buf too small(0) or overflows(-1): ", c)
	}

	fSummaryBytes = fSummaryBytes[c:]
	if err = proto.Unmarshal(fSummaryBytes, fileSummary); err != nil {
		log.Fatal(err)
	}

	for _, value := range fileSummary.GetSections() {
		log.Printf("Section Name: %s:%v", value.GetName(), value)
		sectionMap[value.GetName()] = value
	}
	return sectionMap
}

func ParseInodeSection(info *FileSummary_Section, imageFile *os.File) (map[InodeId]INode, EntityCount) {
	var (
		inodeSectionBytes        = make([]byte, info.GetLength())
		nameIdMap                = make(map[InodeId]INode)
		files             uint32 = 0
		dirs              uint32 = 0
		symlinks          uint32 = 0
		inodeType         INodeSection_INode_Type
	)
	_, err := imageFile.ReadAt(inodeSectionBytes, int64(info.GetOffset()))
	logIfErr(err)

	i, c := binary.Uvarint(inodeSectionBytes)
	if c <= 0 {
		log.Fatal("buf too small(0) or overflows(-1): ", c)
	}
	newPos := uint64(c) + i
	tmpBuf := inodeSectionBytes[c:newPos]

	inodeSection := &INodeSection{}
	if err = proto.Unmarshal(tmpBuf, inodeSection); err != nil {
		log.Fatal(err)
	}
	totalInodes := inodeSection.GetNumInodes()

	for a := uint64(0); a < totalInodes; a++ {
		inodeSectionBytes = inodeSectionBytes[newPos:]
		i, c = binary.Uvarint(inodeSectionBytes)
		if c <= 0 {
			log.Fatal("buf too small(0) or overflows(-1): ", c, a)
		}
		newPos = uint64(c) + i
		tmpBuf = inodeSectionBytes[c:newPos]
		inode := &INodeSection_INode{}
		if err = proto.Unmarshal(tmpBuf, inode); err != nil {
			log.Fatal(err)
		}
		inodeType = inode.GetType()
		id := InodeId(inode.GetId())
		if inodeType == 1 {
			nameIdMap[id] = INode{inode.GetName(), id, FILE}
			files++
		} else if inodeType == 2 {
			nameIdMap[id] = INode{inode.GetName(), id, DIRECTORY}
			dirs++
		} else {
			nameIdMap[id] = INode{inode.GetName(), id, SYMLINK}
			symlinks++
		}
	}
	entityCount := EntityCount{Files: files, Directories: dirs, Symlinks: symlinks}
	return nameIdMap, entityCount
}

func ParseInodeDirectorySection(info *FileSummary_Section, imageFile *os.File) map[ParentId][]uint64 {
	var (
		parChildrenMap = make(map[ParentId][]uint64)
	)
	startPos := int64(info.GetOffset())
	length := info.GetLength()
	dirSectionBytes := make([]byte, length)
	_, err := imageFile.ReadAt(dirSectionBytes, startPos)
	logIfErr(err)
	dirEntry := &INodeDirectorySection_DirEntry{}
	for a := length; a > 0; {
		i, c := binary.Uvarint(dirSectionBytes)
		if c <= 0 {
			log.Fatal("buf too small(0) or overflows(-1)")
		}
		newPos := uint64(c) + i
		tmpBuf := dirSectionBytes[c:newPos]
		if err = proto.Unmarshal(tmpBuf, dirEntry); err != nil {
			log.Fatal(err)
		}
		parent := ParentId(dirEntry.GetParent())
		children := dirEntry.GetChildren()
		parChildrenMap[parent] = children
		a -= newPos
		dirSectionBytes = dirSectionBytes[newPos:]
	}
	return parChildrenMap
}
