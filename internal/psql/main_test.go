package psql_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/adrianbrad/psqldocker"
	"github.com/adrianbrad/psqltest"
)

func TestMain(m *testing.M) {
	// psql connection parameters.
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	// run a new psql docker container.
	c, err := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
		psqldocker.WithSQL(`
		CREATE TABLE users(
			user_id UUID PRIMARY KEY,
			email VARCHAR NOT NULL
		);
		`,
		),
	)
	if err != nil {
		log.Fatalf("err while creating new psql container: %s", err)
	}

	defer func() {
		// close the psql container
		err = c.Close()
		if err != nil {
			log.Printf("err while tearing down db container: %s", err)
		}
	}()

	// compose the psql dsn.
	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	// register the psql container connection details
	// in order to be able to spawn new database connections
	// in an isolated transaction.
	psqltest.RegisterWithPSQLDriver(dsn, "pgx")

	// run the package tests.
	m.Run()
}
