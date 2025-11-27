package main

import (
	"fmt"
	"log"
	"os"

	"github.com/talk2sohail/logger"
	api "github.com/talk2sohail/logger/api/v1"
)

func main() {
	dir, err := os.MkdirTemp("", "logger-example")
	if err != nil {
		log.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	c := logger.Config{}
	c.Segment.MaxStoreBytes = 1024
	c.Segment.MaxIndexBytes = 1024

	l, err := logger.NewLog(dir, c)
	if err != nil {
		log.Fatalf("failed to create log: %v", err)
	}

	records := []*api.Record{
		{Value: []byte("hello world 1")},
		{Value: []byte("hello world 2")},
		{Value: []byte("hello world 3")},
	}

	for _, record := range records {
		off, err := l.Append(record)
		if err != nil {
			log.Fatalf("failed to append record: %v", err)
		}
		fmt.Printf("appended record at offset: %d\n", off)
	}

	for i := uint64(0); i < uint64(len(records)); i++ {
		read, err := l.Read(i)
		if err != nil {
			log.Fatalf("failed to read record: %v", err)
		}
		fmt.Printf("read record: %s\n", string(read.Value))
	}
}
