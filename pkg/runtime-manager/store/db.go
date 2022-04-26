package store

import (
	bolt "go.etcd.io/bbolt"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
)

const (
	runtimeDbPath = "/dev/runtime-manager/checkpoint.db"
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

	if _, err := os.Stat(filepath.Dir(runtimeDbPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(runtimeDbPath), 0755); err != nil {
			klog.Errorf("fail to mkdir for %v", filepath.Dir(runtimeDbPath), err)
			os.Exit(1)
		}
		klog.Infof("create db dir %v", filepath.Dir(runtimeDbPath))
	}
	if _, err := os.Stat(runtimeDbPath); os.IsNotExist(err) {
		if _, err = os.Create(runtimeDbPath); err != nil {
			klog.Errorf("fail to crate %v %v", runtimeDbPath, err)
			os.Exit(1)
		}
		klog.Infof("create db file %v", runtimeDbPath)
	}

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
