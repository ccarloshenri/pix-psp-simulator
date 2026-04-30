package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type ListCobVController struct {
	proc *processor.ListCobVProcessor
}

func NewListCobVController(proc *processor.ListCobVProcessor) *ListCobVController {
	return &ListCobVController{proc: proc}
}

func (c *ListCobVController) Handle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	req := processor.ListCobVRequest{
		Status:           q.Get("status"),
		DataDeVencimento: q.Get("dataDeVencimento"),
		Inicio:           q.Get("inicio"),
		Fim:              q.Get("fim"),
	}

	resp, err := c.proc.Process(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.CobVs.CobVs)
}
