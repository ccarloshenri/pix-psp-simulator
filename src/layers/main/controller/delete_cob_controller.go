package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type DeleteCobController struct {
	proc *processor.DeleteCobProcessor
}

func NewDeleteCobController(proc *processor.DeleteCobProcessor) *DeleteCobController {
	return &DeleteCobController{proc: proc}
}

func (c *DeleteCobController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	_, err := c.proc.Process(processor.DeleteCobRequest{TxID: txid})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
