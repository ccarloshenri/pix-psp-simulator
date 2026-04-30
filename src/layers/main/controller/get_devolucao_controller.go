package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type GetDevolucaoController struct {
	proc *processor.GetDevolucaoProcessor
}

func NewGetDevolucaoController(proc *processor.GetDevolucaoProcessor) *GetDevolucaoController {
	return &GetDevolucaoController{proc: proc}
}

func (c *GetDevolucaoController) Handle(w http.ResponseWriter, r *http.Request) {
	e2eid := r.PathValue("e2eid")
	devID := r.PathValue("id")

	resp, err := c.proc.Process(processor.GetDevolucaoRequest{E2EID: e2eid, DevID: devID})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Devolucao)
}
