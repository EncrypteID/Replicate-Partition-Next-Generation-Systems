package segmentedlog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSegmentReadWrite(t *testing.T) {
	idxFile, err := ioutil.TempFile("", "0001.index")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(idxFile.Name())

	storeFile, err := ioutil.TempFile("", "segment-store-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(storeFile.Name())

	segment, err := NewSegement(idxFile.Name(), storeFile.Name(), 1, &defaultConfig)
	if err != nil {
		t.Fatal(err)
	}

	messages := []string{
		"hello",
		"test",
		"abc",
	}

	var ids []uint64

	for _, message := range messages {
		id, err := segment.write([]byte(message))
		if err != nil {
			t.Error(err)
		}

		ids = append(ids, id)
	}

	for i, id := range ids {
		data, err := segment.read(id)
		if err != nil {
			t.Error(err)
		}
		if string(data) != messages[i] {
			t.Error("data is not the same")
		}
	}
}

func TestRemoveSegment(t *testing.T) {
	i, _ := ioutil.TempFile("", "0001.index-remove")
	s, _ := ioutil.TempFile("", "0001.store-remove")
	_ = i.Close()
	_ = s.Close()

	seg, _ := NewSegement(i.Name(), s.Name(), 1, &defaultConfig)
	err := seg.remove()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(i.Name())
	if !errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
		t.Error("should delete index file")
	}
	_, err = os.Stat(s.Name())
	if !errors.Is(err, os.ErrNotExist) {
		t.Error("should delete store file")
	}
}
