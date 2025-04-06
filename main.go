package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello")
}

// Writes a new file and just deletes the old one (by renaming the new file to the old file)
// If our update is interrupted e.g. by a crash, we still have gthe old file, and concurrent
// readers won't get half-written data
func SaveData(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt())
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return err
	}
	defer func() {
		fp.Close()
		if err != nil {
			os.Remove(tmp)
		}
	}()

	_, err = fp.Write(data)

	if err != nil {
		return err
	}

	err = fp.Sync()

	if err != nil {
		return err
	}

	// renaming a file replaces it atomically, and a delete is unnecessary
	return os.Rename(tmp, path)
}

func SaveData_Old1(path string, data []byte) error {

	/*
		Problems: Must read all data into memory, change it, then write back; only usable for small data
		Has concurrency problems.
	*/

	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(data)

	if err != nil {
		return err
	}

	return fp.Sync()
}

func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}
