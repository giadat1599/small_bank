package api

import (
	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/giadat1599/small_bank/utils"
)

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(0, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}
