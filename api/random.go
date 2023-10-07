package api

import (
	"testing"

	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/giadat1599/small_bank/utils"
	"github.com/stretchr/testify/require"
)

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(0, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	return
}
