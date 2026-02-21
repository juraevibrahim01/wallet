package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/juaevibrahim01/wallet/pkg/types"
)

var ErrPhoneRegistred = errors.New("Phone already registred")
var ErrAmountMMustBepositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("Account not found")
var ErrNotenoughBalance = errors.New("Notenough balance is zero")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []types.Payment
}

func (s *Service) RegistrationAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistred
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMMustBepositive
	}

	var account *types.Account
	for i := range s.accounts {
		if s.accounts[i].ID == accountID {
			account = s.accounts[i]
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	// Зачисление средств пока не рассматривается как платёж
	account.Balance += amount

	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMMustBepositive
	}

	var account *types.Account
	for i := range s.accounts {
		if s.accounts[i].ID == accountID {
			account = s.accounts[i]
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance >= amount {
		return nil, ErrNotenoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, *payment)
	return payment, nil
}
