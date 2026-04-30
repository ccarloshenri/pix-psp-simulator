package controller

import (
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type GetPixController struct {
	proc *processor.GetPixProcessor
}

func NewGetPixController(proc *processor.GetPixProcessor) *GetPixController {
	return &GetPixController{proc: proc}
}

func (c *GetPixController) Handle(w http.ResponseWriter, r *http.Request) {
	e2eid := r.PathValue("e2eid")

	resp, err := c.proc.Process(processor.GetPixRequest{E2EID: e2eid})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp.Pix)
}
