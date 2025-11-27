# Logger

A simple, persisted, append-only log library for Go.

This library provides a file-based, write-ahead log that is both durable and easy to use. It's designed for scenarios where you need to record a sequence of events or messages in a way that can be reliably replayed.

## Features

- **Append-only:** Records are always appended to the end of the log, ensuring that existing data is immutable.
- **Persisted:** Data is written to disk, so it survives application restarts.
- **Segmented:** The log is broken into segments, which makes it easier to manage and compact.
- **Indexed:** Each segment has a corresponding index file, which allows for fast lookups of records by offset.
- **Configurable:** You can configure the maximum size of segments to suit your needs.

## Getting Started

### Installation

To use this library in your Go project, you can use `go get`:

```sh
go get github.com/talk2sohail/logger
```

### Usage

Here's a simple example of how to use the logger:

```go
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
```

## How it Works

The logger is composed of a few key components:

- **Log:** The main entry point for the library. It manages a set of segments.
- **Segment:** A single log file and its corresponding index file.
- **Store:** The file where the records are stored.
- **Index:** The file that maps record offsets to their position in the store file.

When you append a record, it's written to the active segment's store file, and an entry is added to the index file. When a segment reaches its maximum size, a new segment is created.

## Configuration

You can configure the logger by passing a `Config` struct to the `NewLog` function. The following options are available:

- `Segment.MaxStoreBytes`: The maximum size of a segment's store file in bytes.
- `Segment.MaxIndexBytes`: The maximum size of a segment's index file in bytes.
- `Segment.InitialOffset`: The initial offset of the first segment.

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License.