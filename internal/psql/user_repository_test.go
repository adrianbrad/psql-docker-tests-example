package psql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/adrianbrad/psql-docker-tests-example/internal/psql"
	"github.com/adrianbrad/psqldocker"
	"github.com/adrianbrad/psqltest"
	"github.com/lib/pq"
	"github.com/matryer/is"
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
		psqldocker.WithSql(`
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

	// exit code
	var ret int

	defer func() {
		// close the psql container
		err = c.Close()
		if err != nil {
			log.Printf("err while tearing down db container: %s", err)
		}

		// exit with the code provided by executing m.Run().
		os.Exit(ret)
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
	psqltest.Register(dsn)

	// run the package tests.
	ret = m.Run()
}

func TestUserRepository(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validUUID := "29974652-b51f-4826-baf0-c8bd2f1cf4c9"

	t.Run("CreateUser", func(t *testing.T) {
		t.Parallel()

		t.Run("Success1", func(t *testing.T) {
			t.Parallel()

			i := is.New(t)

			err := addUser(ctx, t, psql.User{
				ID:    validUUID,
				Email: "email",
			})
			i.NoErr(err)
		})

		t.Run("Success2", func(t *testing.T) {
			t.Parallel()

			i := is.New(t)

			err := addUser(ctx, t, psql.User{
				ID:    validUUID,
				Email: "email",
			})
			i.NoErr(err)
		})

		t.Run("InvalidID", func(t *testing.T) {
			t.Parallel()

			err := addUser(ctx, t, psql.User{
				ID:    "id",
				Email: "email",
			})

			assertPsqlErr(
				t,
				err,
				"22P02",
				"invalid input syntax for type uuid: \"id\"",
			)
		})

		t.Run("DuplicatePrimaryKey", func(t *testing.T) {
			t.Parallel()

			i := is.New(t)

			err := addUser(ctx, t, psql.User{
				ID:    validUUID,
				Email: "email",
			})
			i.NoErr(err)

			err = addUser(ctx, t, psql.User{
				ID:    validUUID,
				Email: "email",
			})
			assertPsqlErr(
				t,
				err,
				"23505",
				"duplicate key value violates unique constraint \"users_pkey\"",
			)
		})
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Parallel()

		t.Run("Success", func(t *testing.T) {
			t.Parallel()

			i := is.New(t)

			userRepo := newUserRepo(t)

			user := psql.User{
				ID:    validUUID,
				Email: "email",
			}

			err := userRepo.AddUser(
				ctx,
				user,
			)
			i.NoErr(err)

			retrievedUser, err := userRepo.GetUser(ctx, validUUID)
			i.NoErr(err)

			i.Equal(user, retrievedUser)
		})

		t.Run("NotFound", func(t *testing.T) {
			t.Parallel()

			i := is.New(t)

			userRepo := newUserRepo(t)

			_, err := userRepo.GetUser(ctx, validUUID)

			i.True(errors.Is(err, sql.ErrNoRows))
		})

		t.Run("InvalidID", func(t *testing.T) {
			t.Parallel()

			userRepo := newUserRepo(t)

			_, err := userRepo.GetUser(ctx, "id")

			assertPsqlErr(
				t,
				err,
				"22P02",
				"invalid input syntax for type uuid: \"id\"",
			)
		})
	})
}

func addUser(
	ctx context.Context,
	t *testing.T,
	user psql.User,
) error {
	t.Helper()

	return newUserRepo(t).AddUser(ctx, user)
}

func assertPsqlErr(
	t *testing.T,
	err error,
	code pq.ErrorCode,
	message string,
) {
	t.Helper()

	i := is.New(t)

	i.Helper()

	var pqErr *pq.Error

	i.True(errors.As(err, &pqErr))

	i.Equal(code, pqErr.Code)
	i.Equal(message, pqErr.Message)
}

func newUserRepo(t *testing.T) *psql.UserRepository {
	t.Helper()

	db := psqltest.NewTransactionTestingDB(t)
	userRepo := psql.NewUserRepository(db)

	return userRepo
}
