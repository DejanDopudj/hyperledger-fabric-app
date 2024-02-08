package src

type RegisterUser struct {
	BankID    string `json:"bankId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
}

type CreateBankAccount struct {
	UserID    string   `json:"userId" validate:"required"`
	AccountID string   `json:"accountId" validate:"required"`
	Amount    float64  `json:"amount" validate:"required"`
	Currency  string   `json:"currency" validate:"required"`
	CardList  []string `json:"cardlist"`
}

type ManageFunds struct {
	BankID    string  `json:"bankId" validate:"required"`
	UserID    string  `json:"userId" validate:"required"`
	AccountID string  `json:"accountId" validate:"required"`
	Amount    float64 `json:"amount" validate:"required"`
}

type TransferFunds struct {
	UserID        string  `json:"userId" validate:"required"`
	FromBankID    string  `json:"fromBankId" validate:"required"`
	FromAccountID string  `json:"fromAccountId" validate:"required"`
	ToBankID      string  `json:"toBankId" validate:"required"`
	ToAccountID   string  `json:"toAccountId" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
}
