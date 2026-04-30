package processor

import "pix-psp-simulator/src/layers/main/bo"

type ListCobRequest struct {
	Status string
	Inicio string
	Fim    string
}

type ListCobResponse struct {
	Cobs bo.ListCobOutput
}

type ListCobProcessor struct {
	listCobBO *bo.ListCobBO
}

func NewListCobProcessor(listCobBO *bo.ListCobBO) *ListCobProcessor {
	return &ListCobProcessor{listCobBO: listCobBO}
}

func (p *ListCobProcessor) Process(req ListCobRequest) (*ListCobResponse, error) {
	output, err := p.listCobBO.Execute(bo.ListCobInput{
		Status: req.Status,
		Inicio: req.Inicio,
		Fim:    req.Fim,
	})
	if err != nil {
		return nil, err
	}
	return &ListCobResponse{Cobs: *output}, nil
}
