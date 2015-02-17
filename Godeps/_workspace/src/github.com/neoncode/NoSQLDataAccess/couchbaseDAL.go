package DataAccess

import (
	"github.com/couchbaselabs/go-couchbase"
	"fmt"
)

type couchbaseConfig struct {
	DBLocation string
	PoolName string
	BucketName string
}

func GetCouchbaseDAL(dbLocation string, poolName string, bucketName string) (IDataAccess) { 
	//return couchbaseConfig{dbLocation, poolName, bucketName}
	config := new(couchbaseConfig)
	config.DBLocation = dbLocation
	config.PoolName = poolName
	config.BucketName = bucketName
	return config
}

func (config couchbaseConfig) Get(key string, value interface{}) (err error) {
	fmt.Println(config.DBLocation)
	c, err := couchbase.Connect(config.DBLocation)
	if (err != nil) { return }

	pool, err := c.GetPool(config.PoolName)
	if (err != nil) { return }

	bucket, err := pool.GetBucket(config.BucketName)
	if (err != nil) { return }
	defer bucket.Close()

	err = bucket.Get(key, value)
	if (err != nil) { return }

	return
}

func (config couchbaseConfig) Remove(key string) (err error) {
	c, err := couchbase.Connect(config.DBLocation)
	if (err != nil) { return }

	pool, err := c.GetPool(config.PoolName)
	if (err != nil) { return }

	bucket, err := pool.GetBucket(config.BucketName)
	if (err != nil) { return }
	defer bucket.Close()

	err = bucket.Delete(key)

	return
}

func (config couchbaseConfig) Set(key string, value interface{}) (err error) { //key string, 
	c, err := couchbase.Connect(config.DBLocation)
	if (err != nil) { return }

	pool, err := c.GetPool(config.PoolName)
	if (err != nil) { return }

	bucket, err := pool.GetBucket(config.BucketName)
	if (err != nil) { return }
	defer bucket.Close()

	added, err := bucket.Add(key, 0, value)
	if (err != nil) { return }

	if !added {
		//I should modify the value for this key
		err = bucket.Set(key, 0, value)
	}
	return
}
