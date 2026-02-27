package types

// Money представляет собой денежную сумму в минимальных единицах (центы, дирамы и т.д.)
type Money int64

// PaymentCategory представляет собой категорию, в которой был совершён платёж (авто, аптеки, рестораны и т.д.)
type PaymentCategory string

// PaymentStatus предоставляет собой статус платежа.
type PaymentStatus string

// Предопределённые статусы платежей.
const (
	PaymentStatusOk         PaymentStatus = "active"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "IN_PROGRESS"
)

// Payment Представляет информацию о платеже.
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Balance   Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Phone string

// Account представляет информацию о счёте пользователя.
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Error string

func (e Error) Error() string {
	return string(e)
}
