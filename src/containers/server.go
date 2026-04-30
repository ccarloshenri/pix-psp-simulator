package containers

import (
	"fmt"
	"net/http"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/controller"
	"pix-psp-simulator/src/layers/main/implementations/memory"
	uuidgen "pix-psp-simulator/src/layers/main/implementations/uuid"
	"pix-psp-simulator/src/layers/main/processor"
)

type Server struct {
	cfg Config
	mux *http.ServeMux
}

func NewServer(cfg Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	// Repositories
	cobRepo := memory.NewCobRepository()
	cobvRepo := memory.NewCobVRepository()
	pixRepo := memory.NewPixRepository()

	// ID generator
	gen := uuidgen.NewIDGenerator()

	// BOs
	createCobBO := bo.NewCreateCobBO(cobRepo, gen)
	getCobBO := bo.NewGetCobBO(cobRepo)
	updateCobBO := bo.NewUpdateCobBO(cobRepo)
	deleteCobBO := bo.NewDeleteCobBO(cobRepo)
	createCobVBO := bo.NewCreateCobVBO(cobvRepo, gen)
	getCobVBO := bo.NewGetCobVBO(cobvRepo)
	updateCobVBO := bo.NewUpdateCobVBO(cobvRepo)
	simulatePayBO := bo.NewSimulatePaymentBO(cobRepo, cobvRepo, pixRepo, gen)
	getPixBO := bo.NewGetPixBO(pixRepo)
	listPixBO := bo.NewListPixBO(pixRepo)
	createDevBO := bo.NewCreateDevolucaoBO(pixRepo, gen)
	getDevBO := bo.NewGetDevolucaoBO(pixRepo)

	// Processors
	createCobProc := processor.NewCreateCobProcessor(createCobBO)
	getCobProc := processor.NewGetCobProcessor(getCobBO)
	updateCobProc := processor.NewUpdateCobProcessor(updateCobBO)
	deleteCobProc := processor.NewDeleteCobProcessor(deleteCobBO)
	createCobVProc := processor.NewCreateCobVProcessor(createCobVBO)
	getCobVProc := processor.NewGetCobVProcessor(getCobVBO)
	updateCobVProc := processor.NewUpdateCobVProcessor(updateCobVBO)
	simPayProc := processor.NewSimulatePaymentProcessor(simulatePayBO)
	getPixProc := processor.NewGetPixProcessor(getPixBO)
	listPixProc := processor.NewListPixProcessor(listPixBO)
	createDevProc := processor.NewCreateDevolucaoProcessor(createDevBO)
	getDevProc := processor.NewGetDevolucaoProcessor(getDevBO)

	// Controllers
	createCobCtrl := controller.NewCreateCobController(createCobProc)
	getCobCtrl := controller.NewGetCobController(getCobProc)
	updateCobCtrl := controller.NewUpdateCobController(updateCobProc)
	deleteCobCtrl := controller.NewDeleteCobController(deleteCobProc)
	createCobVCtrl := controller.NewCreateCobVController(createCobVProc)
	getCobVCtrl := controller.NewGetCobVController(getCobVProc)
	updateCobVCtrl := controller.NewUpdateCobVController(updateCobVProc)
	simPayCtrl := controller.NewSimulatePaymentController(simPayProc)
	getPixCtrl := controller.NewGetPixController(getPixProc)
	listPixCtrl := controller.NewListPixController(listPixProc)
	createDevCtrl := controller.NewCreateDevolucaoController(createDevProc)
	getDevCtrl := controller.NewGetDevolucaoController(getDevProc)

	// Routes
	s.mux.HandleFunc("POST /cob", createCobCtrl.Handle)
	s.mux.HandleFunc("PUT /cob/{txid}", createCobCtrl.Handle)
	s.mux.HandleFunc("GET /cob/{txid}", getCobCtrl.Handle)
	s.mux.HandleFunc("PATCH /cob/{txid}", updateCobCtrl.Handle)
	s.mux.HandleFunc("DELETE /cob/{txid}", deleteCobCtrl.Handle)

	s.mux.HandleFunc("PUT /cobv/{txid}", createCobVCtrl.Handle)
	s.mux.HandleFunc("GET /cobv/{txid}", getCobVCtrl.Handle)
	s.mux.HandleFunc("PATCH /cobv/{txid}", updateCobVCtrl.Handle)

	s.mux.HandleFunc("POST /pix/simulate", simPayCtrl.Handle)
	s.mux.HandleFunc("GET /pix/{e2eid}", getPixCtrl.Handle)
	s.mux.HandleFunc("GET /pix", listPixCtrl.Handle)

	s.mux.HandleFunc("PUT /pix/{e2eid}/devolucao/{id}", createDevCtrl.Handle)
	s.mux.HandleFunc("GET /pix/{e2eid}/devolucao/{id}", getDevCtrl.Handle)

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	fmt.Printf("PIX PSP Simulator rodando em %s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
