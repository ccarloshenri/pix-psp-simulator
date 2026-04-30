package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type SimulatePaymentRequest struct {
	TxID        string `json:"txid"`
	Valor       string `json:"valor"`
	Infopagador string `json:"infopagador,omitempty"`
}

type SimulatePaymentResponse struct {
	Pix models.Pix
}

type SimulatePaymentProcessor struct {
	simulatePaymentBO *bo.SimulatePaymentBO
}

func NewSimulatePaymentProcessor(simulatePaymentBO *bo.SimulatePaymentBO) *SimulatePaymentProcessor {
	return &SimulatePaymentProcessor{simulatePaymentBO: simulatePaymentBO}
}

func (p *SimulatePaymentProcessor) Process(req SimulatePaymentRequest) (*SimulatePaymentResponse, error) {
	if req.TxID == "" {
		return nil, fmt.Errorf("txid é obrigatório")
	}
	if req.Valor == "" {
		return nil, fmt.Errorf("valor é obrigatório")
	}

	output, err := p.simulatePaymentBO.Execute(bo.SimulatePaymentInput{
		TxID:        req.TxID,
		Valor:       req.Valor,
		Infopagador: req.Infopagador,
	})
	if err != nil {
		return nil, err
	}
	return &SimulatePaymentResponse{Pix: output.Pix}, nil
}
