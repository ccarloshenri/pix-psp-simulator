package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type ListPixController struct {
	proc *processor.ListPixProcessor
}

func NewListPixController(proc *processor.ListPixProcessor) *ListPixController {
	return &ListPixController{proc: proc}
}

func (c *ListPixController) Handle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	req := processor.ListPixRequest{
		TxID:  query.Get("txid"),
		Inicio: query.Get("inicio"),
		Fim:   query.Get("fim"),
	}

	resp, err := c.proc.Process(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Pix)
}
