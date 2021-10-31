package main

import (
	"testing"
)


func TestWallet(t *testing.T) {

	assertBalance := func(t testing.TB, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	assertError := func(t testing.TB, got error, want string) {
		t.Helper()
    	if got == nil {
	        t.Fatal("didn't get an error but wanted one")
	    }

	    if got.Error() != want {
        	t.Errorf("got %q, want %q", got, want)
    	}
	}

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}

        wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}

        wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
    })

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
    	startingBalance := Bitcoin(20)
    	wallet := Wallet{startingBalance}
    	err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, "balance is not enough")
    	assertBalance(t, wallet, startingBalance)
})
}
