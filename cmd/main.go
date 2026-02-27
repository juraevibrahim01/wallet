package main

import (
	"fmt"
	"github.com/juaevibrahim01/wallet/pkg/wallet"
)

func main() {
	acсount_1 := wallet.Service{}

	account, err := acсount_1.RegistrationAccount("+992501012048")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = acсount_1.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMMustBepositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}
		return
	}

	fmt.Println(account.Balance) // 10
}
