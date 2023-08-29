package main

import (
	"fmt"
	kvgo "kvgo"
)

func main() {
	opts := kvgo.DefaultOptions
	opts.DirPath = "/tmp/kvgo"
	db, err := kvgo.Open(opts)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("name"), []byte("kvgo"))
	if err != nil {
		panic(err)
	}
	val, err := db.Get([]byte("name"))
	if err != nil {
		panic(err)
	}
	fmt.Println("val = ", string(val))

	err = db.Delete([]byte("name"))
	if err != nil {
		panic(err)
	}
}
