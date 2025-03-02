package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileMgr struct {
	dbDirectory string
	blockSize   int
	isNew       bool
	openFiles   map[string]*os.File
	mu          sync.Mutex
}

func NewFileMgr(dbDirectory string, blockSize int) (*FileMgr, error) {
	mgr := &FileMgr{
		dbDirectory: dbDirectory,
		blockSize:   blockSize,
		openFiles:   make(map[string]*os.File),
	}

	info, err := os.Stat(dbDirectory)

	if os.IsNotExist(err) {
		err = os.Mkdir(dbDirectory, 0755)
		if err != nil {
			return nil, err
		}
		mgr.isNew = true
	} else if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dbDirectory)
	} else if err != nil {
		return nil, err
	}

	if !mgr.isNew {
		files, err := os.ReadDir(mgr.dbDirectory)
		if err != nil {
			return nil, fmt.Errorf("could not read directory: %s:%v", mgr.dbDirectory, err)
		}

		for _, file := range files {
			if len(file.Name()) >= 4 && file.Name()[:4] == "temp" {
				path := filepath.Join(mgr.dbDirectory, file.Name())
				err = os.Remove(path)
				if err != nil {
					return nil, fmt.Errorf("could not remove temp file: %s:%v", path, err)
				}
			}
		}
	}

	return mgr, nil
}

func (mgr *FileMgr) Read(blk *Block, page *Page) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.getFile(blk.filename)
	if err != nil {
		return err
	}

	_, err = file.Seek(int64(blk.blknum)*int64(mgr.blockSize), io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek to block %d in %s: %v", blk.blknum, blk.filename, err)
	}

	_, err = file.Read(page.Buffer())
	if err != nil {
		return fmt.Errorf("could not read block %d in %s: %v", blk.blknum, blk.filename, err)
	}
	return nil
}

func (mg *FileMgr) Write(blk *Block, page *Page) error {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	file, err := mg.getFile(blk.filename)
	if err != nil {
		return err
	}

	_, err = file.Seek(int64(blk.blknum)*int64(mg.blockSize), io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek to block %d in %s: %v", blk.blknum, blk.filename, err)
	}

	_, err = file.Write(page.Buffer())
	if err != nil {
		return fmt.Errorf("could not write block %d in %s: %v", blk.blknum, blk.filename, err)
	}
	return nil
}

func (mg *FileMgr) Append(filename string) (*Block, error) {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	newBlockNum, err := mg.Size(filename)
	if err != nil {
		return nil, err
	}

	offset := newBlockNum * mg.blockSize
	file, err := mg.getFile(filename)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("could not seek to end of %s: %v", filename, err)
	}

	blk := &Block{filename: filename, blknum: newBlockNum}

	emptyBuffer := make([]byte, mg.blockSize)

	_, err = file.Write(emptyBuffer)
	if err != nil {
		return nil, fmt.Errorf("could not write to %s: %v", filename, err)
	}

	return blk, nil
}

func (mgr *FileMgr) Size(filename string) (int, error) {
	file, err := mgr.getFile(filename)
	if err != nil {
		return 0, fmt.Errorf("can not access %s: %v", filename, err)
	}
	info, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("can not stat %s: %v", filename, err)
	}
	return int(info.Size() / int64(mgr.blockSize)), nil
}

func (mgr *FileMgr) getFile(filename string) (*os.File, error) {
	file, exists := mgr.openFiles[filename]
	if exists {
		return file, nil
	}
	filePath := filepath.Join(mgr.dbDirectory, filename)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	mgr.openFiles[filename] = file
	return file, nil
}
