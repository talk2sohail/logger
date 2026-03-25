package logger

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/talk2sohail/logger/api/v1"
)

func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, log *Log,
	){
		"append and read a record succeeds": testAppendRead,
		"offset out of range error":         testOutOfRangeErr,
		"init with existing segments":       testInitExisting,
		"reader":                            testReader,
		"truncate":                          testTruncate,
		"read across multiple segments":     testReadMultiSegment,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := Config{}
			c.Segment.MaxStoreBytes = 32
			log, err := NewLog(dir, c)
			require.NoError(t, err)

			fn(t, log)
		})
	}
}

func testAppendRead(t *testing.T, log *Log) {
	append := &api.Record{
		Value: []byte("hello world"),
	}
	off, err := log.Append(append)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, append.Value, read.Value)
}

func testOutOfRangeErr(t *testing.T, log *Log) {
	read, err := log.Read(1)
	require.Nil(t, read)
	require.Error(t, err)
}

func testInitExisting(t *testing.T, o *Log) {
	append := &api.Record{
		Value: []byte("hello world"),
	}
	for i := 0; i < 3; i++ {
		_, err := o.Append(append)
		require.NoError(t, err)
	}
	require.NoError(t, o.Close())

	off, err := o.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = o.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)

	n, err := NewLog(o.Dir, o.Config)
	require.NoError(t, err)

	off, err = n.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = n.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)
}

func testReader(t *testing.T, log *Log) {
	append := &api.Record{
		Value: []byte("hello world"),
	}
	off, err := log.Append(append)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	reader := log.Reader()
	b, err := io.ReadAll(reader)
	require.NoError(t, err)

	read := &api.Record{}
	err = read.Unmarshal(bytes.NewReader(b[8:]))
	require.NoError(t, err)
	require.Equal(t, append.Value, read.Value)
	require.Equal(t, uint64(0), read.Offset)
}

func testTruncate(t *testing.T, log *Log) {
	append := &api.Record{
		Value: []byte("hello world"),
	}
	for i := 0; i < 3; i++ {
		_, err := log.Append(append)
		require.NoError(t, err)
	}

	err := log.Truncate(1)
	require.NoError(t, err)

	_, err = log.Read(0)
	require.Error(t, err)
}

func testReadMultiSegment(t *testing.T, log *Log) {
	// MaxStoreBytes is 32, so each segment holds ~1 record.
	// Append enough records to span multiple segments.
	records := []*api.Record{
		{Value: []byte("first")},
		{Value: []byte("second")},
		{Value: []byte("third")},
		{Value: []byte("fourth")},
		{Value: []byte("fifth")},
	}

	for i, rec := range records {
		off, err := log.Append(rec)
		require.NoError(t, err)
		require.Equal(t, uint64(i), off)
	}

	require.Greater(t, len(log.segments), 1,
		"expected multiple segments to test binary search")

	for i, want := range records {
		got, err := log.Read(uint64(i))
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value,
			"mismatch at offset %d", i)
	}

	lowest, err := log.LowestOffset()
	require.NoError(t, err)
	got, err := log.Read(lowest)
	require.NoError(t, err)
	require.Equal(t, records[0].Value, got.Value)

	highest, err := log.HighestOffset()
	require.NoError(t, err)
	got, err = log.Read(highest)
	require.NoError(t, err)
	require.Equal(t, records[len(records)-1].Value, got.Value)

	_, err = log.Read(highest + 1)
	require.Error(t, err)
}
