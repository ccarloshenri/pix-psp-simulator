package controller

import (
	"encoding/json"
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type UpdateCobController struct {
	proc *processor.UpdateCobProcessor
}

func NewUpdateCobController(proc *processor.UpdateCobProcessor) *UpdateCobController {
	return &UpdateCobController{proc: proc}
}

func (c *UpdateCobController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	var req processor.UpdateCobRequest
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
