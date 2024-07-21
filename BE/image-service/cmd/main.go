package main

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/ImageStore/user"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cluster := getClusterConfig()
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to open up a session with the Cassandra database! ", err)
	}
	dbUser := user.NewDbUser(rdb, session)
	server(dbUser)
}

func getClusterConfig() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "my_keyspace"
	cluster.Consistency = gocql.Quorum
	return cluster
}
