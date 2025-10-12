package payment

type TransactionStatus string

const (
	Failed  = TransactionStatus("Failed")
	Success = TransactionStatus("Success")
)
