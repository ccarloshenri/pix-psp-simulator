package controller

import (
	"encoding/json"
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type UpdateCobVController struct {
	proc *processor.UpdateCobVProcessor
}

func NewUpdateCobVController(proc *processor.UpdateCobVProcessor) *UpdateCobVController {
	return &UpdateCobVController{proc: proc}
}

func (c *UpdateCobVController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	var req processor.UpdateCobVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body inválido")
		return
	}
	req.TxID = txid

	_, err := c.proc.Process(req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "atualizado"})
}
