package wallet

import (
	"github.com/google/uuid"
	"github.com/juaevibrahim01/wallet/pkg/types"
	"reflect"
	"testing"
)

func TestService_FindPaymentByID_succes(t *testing.T) {
	// Создаем сервис
	s := NewTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	// сравнение платежа
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("got: %v want: %v", got, payment)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// Creating service
	s := NewTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Fatalf("cant add account, error = %v", err)
		return
	}

	// попробуем найти несуцествующий платеж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): cant create payment, error = %v", err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Error(err)
		return
	}

}

func TestService_Reject_succes(t *testing.T) {
	// Creat service
	s := NewTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// Пробуем отменить платеж
	payment := payments[0]
	err = s.Rejected(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("got: %v want: %v", savedPayment.Status, types.PaymentStatusFail)
		return
	}

	savedAccount, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("got: %v want: %v", savedAccount.Balance, defaultTestAccount.balance)
		return
	}
}
