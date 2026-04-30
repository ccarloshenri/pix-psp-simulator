package controller

import (
	"encoding/json"
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type CreateDevolucaoController struct {
	proc *processor.CreateDevolucaoProcessor
}

func NewCreateDevolucaoController(proc *processor.CreateDevolucaoProcessor) *CreateDevolucaoController {
	return &CreateDevolucaoController{proc: proc}
}

func (c *CreateDevolucaoController) Handle(w http.ResponseWriter, r *http.Request) {
	e2eid := r.PathValue("e2eid")
	devID := r.PathValue("id")

	var req processor.CreateDevolucaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body inválido")
		return
	}

	resp, err := c.proc.Process(e2eid, devID, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Devolucao)
}
