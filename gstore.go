package globalstore

import (
	"crypto/sha256"
	"encoding/base32"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"go.etcd.io/bbolt"
)

func hashstr(s string) string {
	v := sha256.Sum256([]byte(s))
	return base32.StdEncoding.EncodeToString(v[:])
}

type GlobalStore struct {
	filepath string
	fl       *flock.Flock
}

func Open(id string) *GlobalStore {
	flock := flock.New(filepath.Join(os.TempDir(), "gs_"+hashstr(id)+".data.lock"))
	return &GlobalStore{filepath: filepath.Join(os.TempDir(), "gs_"+hashstr(id)+".data"), fl: flock}
}

func (gs *GlobalStore) Get(key string) ([]byte, error) {
	err := gs.fl.Lock()
	if err != nil {
		return nil, err
	}
	defer gs.fl.Unlock()

	db, err := bbolt.Open(gs.filepath, 0666, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var value []byte
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		if b == nil {
			return nil
		}
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (gs *GlobalStore) Set(key string, value []byte) error {
	err := gs.fl.Lock()
	if err != nil {
		return err
	}
	defer gs.fl.Unlock()

	db, err := bbolt.Open(gs.filepath, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("data"))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
	if err != nil {
		return err
	}

	return nil
}

var DefaultGlobalStore = Open("default")

func G(key string) ([]byte, error) {
	return DefaultGlobalStore.Get(key)
}

func S(key string, value []byte) error {
	return DefaultGlobalStore.Set(key, value)
}
