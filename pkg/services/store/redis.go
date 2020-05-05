// +build redis

package store

import (
	"github.com/juju/errors"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/redis"

	"github.com/mgorozii/webhookdelia/pkg/config"
)

func newInternalStore(conf config.Conf) (gokv.Store, error) {
	options := redis.Options{
		Address:  conf.StoreOptions["store_redis_address"].(string),
		Password: conf.StoreOptions["store_redis_password"].(string),
		DB:       conf.StoreOptions["store_redis_db"].(int),
	}

	store, err := redis.NewClient(options)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return store, nil
}
