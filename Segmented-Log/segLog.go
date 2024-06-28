package segmentedlog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/edsrzf/mmap-go"
)

type segmentedlog struct {
	index      *index
	store      *store
	SegementID string
}

// index will store mapping between recordID and recordOffset
// it will maintain it in memory and in index file
type index struct {
	mm      mmap.MMap
	idxFile *os.File
	maxSize uint64
	size    uint64
	id      uint64
	startID uint64
}

// store defines a storage abstraction for the log
// log is append only file
type store struct {
	file    *os.File
	size    uint64
	maxSize uint64
}

func NewSegement(indexFile string, storeFile string, startID uint64, cfg *Config) (*segmentedlog, error) {
	index, err := newIndex(indexFile, cfg, startID)
	if err != nil {
		return err
	}
	store, err := newStore(storeFile, cfg)
	if err != nil {
		return err
	}

	sp := strings.Split(filepath.Base(indexFile), "*")

	return &segmentedlog{
		index:      index,
		store:      store,
		SegementID: sp[0],
	}, nil
}

func (s *segmentedlog) read(id uint64) ([]byte, error) {
	offset, err := s.index.read(id)
	if err != nil {
		return nil, err
	}

	data, err := s.store.read(offset)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *segmentedlog) write(data []byte) (uint64, error) {
	offset, err := s.store.write(data)
	if err != nil {
		return 0, err
	}

	id, err := s.index.write(offset)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *segment) close() error {
	err := s.idx.close()
	if err != nil {
		return err
	}

	err = s.store.close()
	if err != nil {
		return err
	}

	return nil
}

func (s *segmentedlog) remove() error {
	err := s.index.remove()
	if err != nil {
		return err
	}

	return s.store.remove()
}