package processor

import "pix-psp-simulator/src/layers/main/bo"

type ListCobVRequest struct {
	Status           string
	DataDeVencimento string
	Inicio           string
	Fim              string
}

type ListCobVResponse struct {
	CobVs bo.ListCobVOutput
}

type ListCobVProcessor struct {
	listCobVBO *bo.ListCobVBO
}

func NewListCobVProcessor(listCobVBO *bo.ListCobVBO) *ListCobVProcessor {
	return &ListCobVProcessor{listCobVBO: listCobVBO}
}

func (p *ListCobVProcessor) Process(req ListCobVRequest) (*ListCobVResponse, error) {
	output, err := p.listCobVBO.Execute(bo.ListCobVInput{
		Status:           req.Status,
		DataDeVencimento: req.DataDeVencimento,
		Inicio:           req.Inicio,
		Fim:              req.Fim,
	})
	if err != nil {
		return nil, err
	}
	return &ListCobVResponse{CobVs: *output}, nil
}
