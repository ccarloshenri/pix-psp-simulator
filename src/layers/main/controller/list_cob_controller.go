package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type ListCobController struct {
	proc *processor.ListCobProcessor
}

func NewListCobController(proc *processor.ListCobProcessor) *ListCobController {
	return &ListCobController{proc: proc}
}

func (c *ListCobController) Handle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	req := processor.ListCobRequest{
		Status: q.Get("status"),
		Inicio: q.Get("inicio"),
		Fim:    q.Get("fim"),
	}

	resp, err := c.proc.Process(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Cobs.Cobs)
}
