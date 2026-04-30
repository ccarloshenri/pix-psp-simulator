package controller

import (
	"encoding/json"
	"net/http"

	"pix-psp-simulator/src/layers/main/processor"
)

type SimulatePaymentController struct {
	proc *processor.SimulatePaymentProcessor
}

func NewSimulatePaymentController(proc *processor.SimulatePaymentProcessor) *SimulatePaymentController {
	return &SimulatePaymentController{proc: proc}
}

func (c *SimulatePaymentController) Handle(w http.ResponseWriter, r *http.Request) {
	var req processor.SimulatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body inválido")
		return
	}

	resp, err := c.proc.Process(req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, resp)
}
