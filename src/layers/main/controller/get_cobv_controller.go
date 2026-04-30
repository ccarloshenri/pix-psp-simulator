package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type GetCobVController struct {
	proc *processor.GetCobVProcessor
}

func NewGetCobVController(proc *processor.GetCobVProcessor) *GetCobVController {
	return &GetCobVController{proc: proc}
}

func (c *GetCobVController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	resp, err := c.proc.Process(processor.GetCobVRequest{TxID: txid})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.CobV)
}
