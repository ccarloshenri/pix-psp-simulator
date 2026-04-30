package interfaces

type PaymentJob struct {
	E2EID       string
	TxID        string
	Valor       string
	Infopagador string
}

type PaymentQueue interface {
	Enqueue(job PaymentJob)
	Jobs() <-chan PaymentJob
}
