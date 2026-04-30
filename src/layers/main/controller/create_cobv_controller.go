package controller

import (
	"encoding/json"
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type CreateCobVController struct {
	proc *processor.CreateCobVProcessor
}

func NewCreateCobVController(proc *processor.CreateCobVProcessor) *CreateCobVController {
	return &CreateCobVController{proc: proc}
}

func (c *CreateCobVController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	var req processor.CreateCobVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body inválido")
		return
	}

	resp, err := c.proc.Process(txid, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, resp.CobV)
}
