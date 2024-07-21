package user

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
)

type DbUser struct {
	RClient        *redis.Client
	CassConnection *gocql.Session
}

func NewDbUser(rClient *redis.Client, cSess *gocql.Session) *DbUser {
	return &DbUser{
		RClient:        rClient,
		CassConnection: cSess,
	}
}

func (r *DbUser) UserExists(key string) (error, string) {
	val, err := r.RClient.Get(key).Result()
	return err, val
}

func (r *DbUser) GetUser(key string) (string, error) {
	user, err := r.RClient.Get(key).Result()
	return user, err
}
func (r *DbUser) SetUser(key string, val string) error {
	_, err := r.RClient.Set(key, val, 2*time.Minute).Result()
	return err
}
