// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"errors"
	"time"

	gocb "github.com/couchbase/gocb/v2"
)

const (
	// Duration to wait until memd connections have been established with
	// the server and are ready.
	timeout = 30 * time.Second
)

var (
	// ErrDocumentNotFound is returned when the doc does not exist in the DB.
	ErrDocumentNotFound = errors.New("document not found")
)

// DB represents the database connection.
type DB struct {
	Bucket     *gocb.Bucket
	Cluster    *gocb.Cluster
	Collection *gocb.Collection
}

// Config represents the database config.
type Config struct {
	// the data source name (DSN) for connecting to the database.
	Server string `mapstructure:"server"`
	// Username used to access the db.
	Username string `mapstructure:"username"`
	// Password used to access the db.
	Password string `mapstructure:"password"`
	// Name of the couchbase bucket.
	BucketName string `mapstructure:"bucket_name"`
}

// Open opens a connection to the database.
func Open(server, username, password, bucketName string) (DB, error) {

	// Get a couchbase cluster instance.
	cluster, err := gocb.Connect(
		server,
		gocb.ClusterOptions{
			Username: username,
			Password: password,
		})
	if err != nil {
		return DB{}, err
	}

	// Get a bucket reference.
	bucket := cluster.Bucket(bucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(timeout, nil)
	if err != nil {
		return DB{}, err
	}

	// Get a collection reference.
	collection := bucket.DefaultCollection()

	// Create primary index.
	mgr := cluster.QueryIndexes()
	err = mgr.CreatePrimaryIndex(bucketName,
		&gocb.CreatePrimaryQueryIndexOptions{IgnoreIfExists: true})
	if err != nil {
		return DB{}, err
	}

	return DB{
		Bucket:     bucket,
		Cluster:    cluster,
		Collection: collection,
	}, nil
}

// Query executes a N1QL query.
func (db *DB) Query(ctx context.Context, statement string,
	args map[string]interface{}) (*gocb.QueryResult, error) {

	results, err := db.Cluster.Query(statement, &gocb.QueryOptions{
		NamedParameters: args, Adhoc: true})
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Get retrieves the document using its key.
func (db *DB) Get(ctx context.Context, key string, model interface{}) error {

	// Performs a fetch operation against the collection.
	getResult, err := db.Collection.Get(key, &gocb.GetOptions{})

	if errors.Is(err, gocb.ErrDocumentNotFound) {
		return ErrDocumentNotFound
	}
	if err != nil {
		return err
	}

	// Assigns the value of the result into the valuePtr using default decoding.
	err = getResult.Content(&model)
	if err != nil {
		return err
	}

	return nil
}

// Create saves a new document into the collection.
func (db *DB) Create(ctx context.Context, key string, val interface{}) error {
	_, err := db.Collection.Insert(key, val, &gocb.InsertOptions{})
	return err
}

// Update updates a document in the collection.
func (db *DB) Update(ctx context.Context, key string, path string,
	val interface{}) error {

	// When `path` is not empty, we performs a sub document in the collection.
	// Sub documents operations may be quicker and more network-efficient than
	// full-document operations.
	if len(path) > 0 {
		mops := []gocb.MutateInSpec{
			gocb.UpsertSpec(path, val, &gocb.UpsertSpecOptions{CreatePath: true}),
		}
		_, err := db.Collection.MutateIn(key, mops,
			&gocb.MutateInOptions{Timeout: 10050 * time.Millisecond})
		return err
	}

	_, err := db.Collection.Replace(key, val, &gocb.ReplaceOptions{})
	return err
}

// Delete removes a document from the collection.
func (db *DB) Delete(ctx context.Context, key string) error {
	_, err := db.Collection.Remove(key, &gocb.RemoveOptions{})
	return err
}

// Count retrieves the total number of documents.
func (db *DB) Count(ctx context.Context, docType string,
	val interface{}) error {

	//val = nil

	params := make(map[string]interface{}, 2)
	params["bucketName"] = db.Bucket.Name()
	params["docType"] = docType

	statement := `SELECT COUNT(*) FROM $bucketname WHERE type=$docType`
	results, err := db.Query(ctx, statement, params)
	if err != nil {
		return err
	}

	var row int
	err = results.One(&row)
	if err != nil {
		return err
	}

	return nil
}

// Lookup query the document for certain path(s); these path(s) are then returned.
func (db *DB) Lookup(ctx context.Context, key string, path string,
	val interface{}) error {

	ops := []gocb.LookupInSpec{gocb.GetSpec(path, &gocb.GetSpecOptions{})}
	getResult, err := db.Collection.LookupIn(key, ops, &gocb.LookupInOptions{})
	if err != nil {
		return err
	}

	return getResult.ContentAt(0, &val)
}
