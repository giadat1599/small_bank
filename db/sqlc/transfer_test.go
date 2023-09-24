package db

import (
	"context"
	"testing"
	"time"

	"github.com/giadat1599/small_bank/utils"
	"github.com/stretchr/testify/require"
)

func create2RandomAccounts(t *testing.T) (Account, Account) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	return account1, account2
}

func createRandomTransfer(t *testing.T, accountFrom Account, accountTo Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	accountFrom, accountTo := create2RandomAccounts(t)
	createRandomTransfer(t, accountFrom, accountTo)
}

func TestGetTransfer(t *testing.T) {
	accountFrom, accountTo := create2RandomAccounts(t)
	transfer1 := createRandomTransfer(t, accountFrom, accountTo)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestGetListTransfers(t *testing.T) {
	accountFrom, accountTo := create2RandomAccounts(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, accountFrom, accountTo)
	}

	arg := ListTransfersParams{
		FromAccountID: accountFrom.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
