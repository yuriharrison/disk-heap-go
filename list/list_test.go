package list

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"testing"

	uuid "github.com/satori/uuid"
)

const (
	file        = "/../internal/test.db"
	float64Size = 8
	uuidSize    = 16
	itemSize    = float64Size + uuidSize
	// 500k items per page
	maxPageSize = 10_000_000
	// 100 items per page
	minPageSize = itemSize * 100
)

type testItem struct {
	idx   int
	id    uuid.UUID
	value float64
}

func getTestData() []testItem {
	return []testItem{
		{0, uuid.Must(uuid.NewV4()), 100},
		{1, uuid.Must(uuid.NewV4()), 200},
		{2, uuid.Must(uuid.NewV4()), 300},
		{3, uuid.Must(uuid.NewV4()), 400},
		{4, uuid.Must(uuid.NewV4()), 500},
	}
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

func serializeItem(item testItem) []byte {
	data := make([]byte, itemSize)
	for i, v := range item.id {
		data[i] = v
	}
	for i, v := range float64bytes(item.value) {
		data[i+uuidSize] = v
	}
	return data
}

func deserialize(data []byte) testItem {
	fmt.Println(data)
	var id [uuidSize]byte
	copy(id[:], data[:uuidSize])
	floatByte := make([]byte, float64Size)
	copy(floatByte, data[uuidSize:])
	fmt.Println(floatByte, data)
	value := float64frombytes(floatByte)
	return testItem{id: id, value: value}
}

func createFile(filePath string) {
	f := openFile(filePath, os.O_CREATE)
	f.Write(make([]byte, itemSize))
	f.Close()
}

func TestList(t *testing.T) {
	data := getTestData()
	wd, _ := os.Getwd()
	filePath := wd + file
	createFile(filePath)
	l := NewList(itemSize, filePath, maxPageSize, minPageSize)
	for _, d := range data {
		l.Append(serializeItem(d))
	}
	for _, d := range data {
		bdt := l.Get(d.idx)
		if bdt == nil || len(bdt) < itemSize {
			t.Error("Invalid value!")
		}
		item := deserialize(bdt)
		if item.value != d.value {
			t.Errorf("Wrong value on index %v value %v != %v", d.idx, d.value, item.value)
		}
	}
}
