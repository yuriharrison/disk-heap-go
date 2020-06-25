package list

import (
	"encoding/binary"
	"log"
	"math"
	"os"

	"github.com/edsrzf/mmap-go"
	uuid "github.com/satori/go.uuid"
)

const (
	fileName    = "/internal/data.db"
	float64Size = 8
	uuidSize    = 16
	nodeSize    = float64Size + uuidSize
	// 500k items per page
	maxPageSize = 10_000_000
	// 100 items per page
	minPageSize = nodeSize * 100
)

type NodeList struct {
	sz         uint32
	dbPath     string
	osPageSize int
}

// NewNodeList Create new NodeList instance
func NewNodeList() NodeList {
	wdir, _ := os.Getwd()
	dbPath := wdir + fileName
	return NodeList{dbPath: dbPath, osPageSize: os.Getpagesize()}
}

func openFile(filePath string, flags int) *os.File {
	f, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, float64Size)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func (l *NodeList) Append(index uint32, id uuid.UUID, value float64) {
	// TODO transform id and value in a Node struct
	data := make([]byte, nodeSize)
	for i, v := range id {
		data[i] = v
	}
	for i, v := range float64bytes(value) {
		data[i] = v
	}
	l.write(index, data)
	l.sz++
}

func (l *NodeList) Set() {

}

func (l *NodeList) Get() {

}

func (l *NodeList) Size() uint32 {
	return l.sz
}

func (l *NodeList) write(index uint32, data []byte) {
	var pageSize, innerPageOffset int
	var offset int64
	fileSize := int(l.sz * nodeSize)
	if fileSize > maxPageSize {
		pageSize = maxPageSize
		page := math.Floor(float64(int(index) * nodeSize / pageSize))
		offset = int64(int(page) * pageSize)
		innerPageOffset = int(index) * nodeSize % pageSize
	} else {
		innerPageOffset = fileSize
		pageSize = fileSize + nodeSize
		if fileSize >= l.osPageSize {
			for pageSize%l.osPageSize != 0 {
				pageSize++
			}
		} else {
			pageSize = l.osPageSize
		}
	}

	f := openFile(l.dbPath, os.O_RDWR)
	m, err := mmap.MapRegion(f, pageSize, mmap.RDWR, 0, offset)
	if err != nil {
		log.Fatal("Error mapping", err)
	}

	for _, b := range data {
		m[innerPageOffset] = b
		innerPageOffset++
	}
}
