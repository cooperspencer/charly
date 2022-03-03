package db

import (
	"charly/types"
	"encoding/json"
	"errors"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
)

type DB struct {
	db *bolt.DB
}

func New(path string, mode os.FileMode) (DB, error) {
	db, err := bolt.Open(path, mode, nil)
	return DB{db}, err
}

func (d DB) InsertRepo(repo types.Repo) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("repos"))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(repo)
		if err != nil {
			return err
		}

		return b.Put([]byte(fmt.Sprintf("%s:%s", repo.Branch, repo.Url)), buf)
	})
}

func (d DB) GetRepo(repo types.Repo) (types.Repo, error) {
	r := types.Repo{}
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("repos"))

		if b == nil {
			return errors.New("repos bucket doesn't exist")
		}

		data := b.Get([]byte(fmt.Sprintf("%s:%s", repo.Branch, repo.Url)))
		err := json.Unmarshal(data, &r)
		return err
	})

	return r, err
}
