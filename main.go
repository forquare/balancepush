package main

import (
	"fmt"
	gocardless "github.com/forquare/balancepush-gocardless"
	"github.com/forquare/balancepush/config"
	"github.com/gregdel/pushover"
	"log"
	"strings"
)

func main() {
	c := config.GetConfig()

	client, err := gocardless.NewGoCardlessClient(c.GoCardless.Credentials.SecretID, c.GoCardless.Credentials.SecretKey)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	l := len(c.GoCardless.Bank.Accounts)
	msg := make([]string, l)
	for i, account := range c.GoCardless.Bank.Accounts {
		account.Balance, account.Currency, account.CurrencySymbol, err = client.GetAccountBalance(account.ID, account.BalanceType)
		if err != nil {
			log.Fatalf("Could not get balance: %v", err)
		}

		m := fmt.Sprintf("%s: %s%.2f", account.Name, account.CurrencySymbol, account.Balance)
		msg[i] = m
	}

	app := pushover.New(c.Pushover.Tokens.App)
	recipient := pushover.NewRecipient(c.Pushover.Tokens.User)
	message := pushover.NewMessage(strings.Join(msg, "\n"))
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(response)
}
