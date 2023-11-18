package db

import "context"

type CreateUserTXParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTXResult struct {
	User User
}

func (store *SQLStore) CreateUserTX(ctx context.Context, arg CreateUserTXParams) (CreateUserTXResult, error) {
	var result CreateUserTXResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)

		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)

	})

	return result, err
}
