package main

import (
	"os"
)

func openFile(filePath string, flags int) *os.File {
	f, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

// func getFilePath() string {
// 	wd, _ := os.Getwd()
// 	return wd + "/internal/data"
// }

// func writeData(data []byte) *os.File {
// 	f := openFile(getFilePath(), os.O_RDWR)
// 	f.Write(data)
// 	return f
// }

// func modify(file *os.File, pageSize int, offset int64) {
// 	// m, err := mmap.Map(f, mmap.RDWR, 0)
// 	m, err := mmap.MapRegion(file, pageSize, mmap.RDWR, 0, offset)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer m.Unmap()

// 	m[1025] = '!'
// 	m[1488] = '&'
// 	m.Flush()
// }

// func assertModification(file *os.File, pageSize int, offset int64) {
// 	m, err := mmap.MapRegion(file, pageSize, mmap.RDONLY, 0, offset)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer m.Unmap()

// 	errMsg := ""
// 	if m[1025] != '!' {
// 		errMsg += "wrong !"
// 	}
// 	if m[1488] != '&' {
// 		errMsg += "wrong &"
// 	}
// 	if errMsg != "" {
// 		panic(errMsg)
// 	}
// }

func main() {
	// page := 2
	// pageSize := 40 * 1024
	// offset := int64((page - 1) * pageSize)
	// var testData []byte
	// for i := 0; i < pageSize*2; i++ {
	// 	testData = append(testData, '#')
	// }
	// f := writeData(testData)
	// defer f.Close()

	// modify(f, pageSize, offset)
	// assertModification(f, pageSize, offset)

	// v, _ := os.Getwd()
	// f := openFile(v+"/internal/data.db", os.O_RDWR)
	// f.Write(make([]byte, os.Getpagesize()))
	// f.Close()
}
