package list

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/edsrzf/mmap-go"
)

type List struct {
	sz          int
	dbPath      string
	osPageSize  int
	itemSize    int
	maxPageSize int
	minPageSize int
}

// NewList Create new List instance
func NewList(itemSize int, dbPath string, maxPageSize, minPageSize int) List {
	return List{
		dbPath:      dbPath,
		osPageSize:  os.Getpagesize(),
		itemSize:    itemSize,
		maxPageSize: maxPageSize,
		minPageSize: minPageSize,
	}
}

func openFile(filePath string, flags int) *os.File {
	f, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func (l *List) validateData(data []byte) {
	if l.osPageSize%len(data) == 0 {
		panic(fmt.Sprintf("Data must be divisor of %v", l.osPageSize))
	}
}

// Append new item
func (l *List) Append(data []byte) int {
	nextIndex := l.sz
	l.validateData(data)
	l.write(nextIndex, data)
	l.sz++
	return nextIndex
}

func (l *List) Set(index int, data []byte) {
	if index > l.sz-1 {
		panic(fmt.Sprintf("Index out of range %v", index))
	}
	l.validateData(data)
}

func (l *List) Get(index int) []byte {
	if index > l.sz-1 {
		return nil
	}
	data := make([]byte, l.itemSize)
	l.read(index, data)
	return data
}

func (l *List) Size() int {
	return l.sz
}

func (l *List) fileOffset(index int) (offset int64, pageSize, innerPageOffset int) {
	fileSize := int(l.sz * l.itemSize)
	innerPageOffset = index * l.itemSize
	if fileSize > l.maxPageSize {
		pageSize = l.maxPageSize
		page := math.Floor(float64(index * l.itemSize / pageSize))
		offset = int64(int(page) * pageSize)
		innerPageOffset = innerPageOffset % pageSize
	} else {
		pageSize = fileSize + l.itemSize
		if fileSize >= l.osPageSize {
			for pageSize%l.osPageSize != 0 {
				pageSize++
			}
		} else {
			pageSize = l.osPageSize
		}
	}
	return
}

func (l *List) getMap(pageSize int, offset int64) mmap.MMap {
	f := openFile(l.dbPath, os.O_RDWR)
	m, err := mmap.MapRegion(f, pageSize, mmap.RDWR, 0, offset)
	if err != nil {
		log.Fatal("Error mapping file:", err)
	}
	return m
}

func (l *List) write(index int, data []byte) {
	offset, pageSize, innerPageOffset := l.fileOffset(index)
	m := l.getMap(pageSize, offset)
	println("\n\n\n", "data", len(data), "m", len(m), "poggset", innerPageOffset, "\n\n\n")
	for _, b := range data {
		m[innerPageOffset] = b
		innerPageOffset++
	}
	err := m.Flush()
	if err != nil {
		log.Fatal("Error flushing data.")
	}
}

func (l *List) read(index int, data []byte) {
	offset, pageSize, innerPageOffset := l.fileOffset(index)
	m := l.getMap(pageSize, offset)
	for i := range data {
		data[i] = m[innerPageOffset]
		innerPageOffset++
	}
}
