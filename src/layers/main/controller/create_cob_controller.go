package controller

import (
	"encoding/json"
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type CreateCobController struct {
	proc *processor.CreateCobProcessor
}

func NewCreateCobController(proc *processor.CreateCobProcessor) *CreateCobController {
	return &CreateCobController{proc: proc}
}

func (c *CreateCobController) Handle(w http.ResponseWriter, r *http.Request) {
	txid := r.PathValue("txid")

	var req processor.CreateCobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body inválido")
		return
	}

	resp, err := c.proc.Process(txid, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, resp.Cob)
}
