package objectstorage

import (
	"github.com/iotaledger/hive.go/parameter"
	"os"
	"sync"

	"github.com/pkg/errors"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
)

var instance *badger.DB

var once sync.Once

func GetBadgerInstance(optionalDirectory ...string) *badger.DB {
	once.Do(func() {
		var directory string
		if len(optionalDirectory) >= 1 {
			directory = optionalDirectory[0]
		} else {
			directory = parameter.DefaultConfig().GetString("objectstorage.directory")
		}

		db, err := createDB(directory)
		if err != nil {
			// errors should cause a panic to avoid singleton deadlocks
			panic(err)
		}
		instance = db
	})

	return instance
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkDir(dir string) error {
	exists, err := exists(dir)
	if err != nil {
		return err
	}

	if !exists {
		return os.Mkdir(dir, 0700)
	}
	return nil
}

func createDB(directory string) (*badger.DB, error) {
	if err := checkDir(directory); err != nil {
		return nil, errors.Wrap(err, "Could not check directory")
	}

	opts := badger.DefaultOptions(directory)
	opts.Logger = nil
	opts.Truncate = true
	opts.TableLoadingMode = options.MemoryMap

	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open new DB")
	}

	return db, nil
}
