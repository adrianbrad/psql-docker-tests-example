package psql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/adrianbrad/psql-docker-tests-example/internal/psql"
	"github.com/adrianbrad/psqltest"
	"github.com/lib/pq"
	"github.com/matryer/is"
)

// func TestUserRepository(t *testing.T) {
// 	t.Parallel()
//
// 	t.Run("CreateUser", func(t *testing.T) {
// 		t.Parallel()
//
// 		t.Run("Success1", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
//
// 		t.Run("Success2", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
//
// 		t.Run("InvalidID", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
//
// 		t.Run("DuplicatePrimaryKey", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
// 	})
//
// 	t.Run("GetUser", func(t *testing.T) {
// 		t.Parallel()
//
// 		t.Run("Success", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
//
// 		t.Run("NotFound", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
//
// 		t.Run("InvalidID", func(t *testing.T) {
// 			t.Parallel()
//
// 			db := psqltest.NewTransactionTestingDB(t)
// 			...
// 		})
// 	})
// }

func addUser(
	ctx context.Context,
	t *testing.T,
	user psql.User,
) error {
	t.Helper()

	return psql.
		NewUserRepository(psqltest.NewTransactionTestingDB(t)).
		AddUser(ctx, user)
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

func TestUserRepository(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validUUID := "29974652-b51f-4826-baf0-c8bd2f1cf4c9"

	t.Run("CreateUser", func(t *testing.T) {
		t.Parallel()

		t.Run("Success1", func(t *testing.T) {
			t.Parallel()

			i := is.New(t)

			fmt.Println(t.Name())

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

			userRepo := psql.
				NewUserRepository(psqltest.NewTransactionTestingDB(t))

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

			userRepo := psql.
				NewUserRepository(psqltest.NewTransactionTestingDB(t))

			_, err := userRepo.GetUser(ctx, validUUID)

			i.True(errors.Is(err, sql.ErrNoRows))
		})

		t.Run("InvalidID", func(t *testing.T) {
			t.Parallel()

			userRepo := psql.
				NewUserRepository(psqltest.NewTransactionTestingDB(t))

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
