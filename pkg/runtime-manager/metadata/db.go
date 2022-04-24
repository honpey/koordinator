package metadata

import (
	bolt "go.etcd.io/bbolt"
	"k8s.io/klog/v2"
	"os"
)

const (
	runtimeDbPath = "/dev/runtime-manager/podmeta.db"
)

var (
	PodBucket       = []byte("podBucket")
	ContainerBucket = []byte("containerBucket")
)

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB() *BoltDB {
	return &BoltDB{}
}

func (m *BoltDB) Init() {
	options := *bolt.DefaultOptions
	options.Timeout = 0

	db, err := bolt.Open(runtimeDbPath, 0644, &options)
	if err != nil {
		klog.Errorf("fail to open bolt db %v", err)
		os.Exit(1)
	}
	m.db = db
	// create 2 bucket for podmeta info store and container metainfo store
	if err := m.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(PodBucket); err != nil {
			klog.Errorf("fail to create bucket %v", err)
		}
		if _, err := tx.CreateBucketIfNotExists(ContainerBucket); err != nil {
			klog.Errorf("fail to create bucket %v", err)
		}

		return nil
	}); err != nil {
	}
}

func (m *BoltDB) Update(namespace, key, val []byte) error {
	if err := m.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(namespace).Put(key, val)
	}); err != nil {
		return err
	}
	return nil
}

func (m *BoltDB) Read(namespace, key []byte) ([]byte, error) {
	var out []byte
	if err := m.db.View(func(tx *bolt.Tx) error {
		out = tx.Bucket(namespace).Get(key)
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
}
