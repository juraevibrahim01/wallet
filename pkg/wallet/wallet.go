package wallet

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/juaevibrahim01/wallet/pkg/types"
)

var ErrPhoneRegistred = errors.New("Phone already registred")
var ErrAmountMMustBepositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("Account not found")
var ErrNotenoughBalance = errors.New("Notenough balance is zero")
var ErrPaymentNotFound = errors.New("Payment not found")

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

func (s *Service) Rejected(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return &payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

type testService struct {
	*Service
}

func NewTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	// Регистрируем там пользовател
	account, err := s.RegistrationAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("cant register account, error = %v", err)
	}

	// Пополняем его счет
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("cant deposit account, error = %v", err)
	}

	return account, nil
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992501012048",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{{amount: 1_000_00, category: "auto"}},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// регистрируем там пользователя
	account, err := s.RegistrationAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register account, error = %v", err)
	}

	// пополняем его счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cant deposit account, error = %v", err)
	}

	// Выполняем платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cant register payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}
