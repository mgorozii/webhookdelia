// +build postgres

package store

import (
	"github.com/juju/errors"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/postgresql"

	"github.com/mgorozii/webhookdelia/pkg/config"
)

func newInternalStore(conf config.Conf) (gokv.Store, error) {
	options := postgresql.Options{
		ConnectionURL:      conf.StoreOptions["store_postgres_connection_url"].(string),
		TableName:          conf.StoreOptions["store_postgres_table_name"].(string),
		MaxOpenConnections: conf.StoreOptions["store_postgres_max_open_connections"].(int),
	}

	store, err := postgresql.NewClient(options)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return store, nil
}
