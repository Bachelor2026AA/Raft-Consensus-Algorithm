package db

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

func BboltDatabase(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()
	return db, nil

}
