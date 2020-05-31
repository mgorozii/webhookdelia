// +build !postgres,!redis

package store

import (
	"github.com/juju/errors"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/file"

	"github.com/mgorozii/webhookdelia/internal/config"
)

func newInternalStore(conf config.Conf) (gokv.Store, error) {
	options := file.Options{
		Directory: conf.StoreOptions["store_file_directory"].(string),
	}

	store, err := file.NewStore(options)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return store, nil
}
