package arangodb

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type service struct {
	conn   driver.Connection
	client driver.Client
	db     driver.Database
}

func NewService(endpoint string, dbname string) (*service, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "password"),
	})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	var db driver.Database
	if ok, _ := client.DatabaseExists(ctx, dbname); !ok {
		db, err = client.CreateDatabase(ctx, dbname, nil)
		if err != nil {
			return nil, err
		}
	} else {
		db, err = client.Database(ctx, dbname)
		if err != nil {
			return nil, err
		}

	}
	return &service{
		conn:   conn,
		client: client,
		db:     db,
	}, nil
}

func (s *service) Collection(ctx context.Context, name string) (driver.Collection, error) {
	if ok, _ := s.db.CollectionExists(ctx, name); !ok {
		return s.db.CreateCollection(ctx, name, nil)
	} else {
		return s.db.Collection(ctx, name)
	}
}
