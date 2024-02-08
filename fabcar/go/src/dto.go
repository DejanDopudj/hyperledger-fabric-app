package src

type RegisterUser struct {
	BankID    string `json:"bankId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
}

type CreateBankAccount struct {
	BankID    string   `json:"bankId" validate:"required"`
	UserID    string   `json:"userId" validate:"required"`
	AccountID string   `json:"accountId" validate:"required"`
	Amount    float64  `json:"amount" validate:"required"`
	Currency  string   `json:"currency" validate:"required"`
	CardList  []string `json:"cardlist"`
}

type Payment struct {
	BankID    string  `json:"bankId" validate:"required"`
	UserID    string  `json:"userId" validate:"required"`
	AccountID string  `json:"accountId" validate:"required"`
	Amount    float64 `json:"amount" validate:"required"`
	Currency  string  `json:"currency" validate:"required"`
}

type Withdrawal struct {
	BankID    string  `json:"bankId" validate:"required"`
	UserID    string  `json:"userId" validate:"required"`
	AccountID string  `json:"accountId" validate:"required"`
	Amount    float64 `json:"amount" validate:"required"`
}

type TransferFunds struct {
	UserID           string  `json:"userId" validate:"required"`
	BankID           string  `json:"fromBankId" validate:"required"`
	FromAccountID    string  `json:"fromAccountId" validate:"required"`
	ToAccountID      string  `json:"toAccountId" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	AcceptConversion bool    `json:"acceptConversion" validate:"required"`
}
