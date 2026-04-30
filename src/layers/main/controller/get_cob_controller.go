package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type GetCobController struct {
	proc *processor.GetCobProcessor
}

func NewGetCobController(proc *processor.GetCobProcessor) *GetCobController {
	return &GetCobController{proc: proc}
}

func (c *GetCobController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	resp, err := c.proc.Process(processor.GetCobRequest{TxID: txid})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Cob)
}
