package containers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/controller"
	channelqueue "pix-psp-simulator/src/layers/main/implementations/channel"
	"pix-psp-simulator/src/layers/main/implementations/memory"
	sqlrepo "pix-psp-simulator/src/layers/main/implementations/sql"
	"pix-psp-simulator/src/layers/main/implementations/worker"
	uuidgen "pix-psp-simulator/src/layers/main/implementations/uuid"
	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/processor"
)

type Server struct {
	cfg    Config
	mux    *http.ServeMux
	worker *worker.PaymentWorker
}

func NewServer(cfg Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}

	cobRepo, cobvRepo, pixRepo := buildRepositories(cfg)

	// ID generator
	gen := uuidgen.NewIDGenerator()

	// Async payment queue and worker
	queue := channelqueue.NewQueue(256)
	s.worker = worker.NewPaymentWorker(queue, cobRepo, cobvRepo, pixRepo)
	go s.worker.Start(context.Background())

	// BOs
	createCobBO := bo.NewCreateCobBO(cobRepo, gen)
	getCobBO := bo.NewGetCobBO(cobRepo)
	listCobBO := bo.NewListCobBO(cobRepo)
	updateCobBO := bo.NewUpdateCobBO(cobRepo)
	deleteCobBO := bo.NewDeleteCobBO(cobRepo)
	createCobVBO := bo.NewCreateCobVBO(cobvRepo, gen)
	getCobVBO := bo.NewGetCobVBO(cobvRepo)
	listCobVBO := bo.NewListCobVBO(cobvRepo)
	updateCobVBO := bo.NewUpdateCobVBO(cobvRepo)
	simulatePayBO := bo.NewSimulatePaymentBO(cobRepo, cobvRepo, gen, queue)
	createDevBO := bo.NewCreateDevolucaoBO(pixRepo, gen)
	getDevBO := bo.NewGetDevolucaoBO(pixRepo)

	// Processors
	createCobProc := processor.NewCreateCobProcessor(createCobBO)
	getCobProc := processor.NewGetCobProcessor(getCobBO)
	listCobProc := processor.NewListCobProcessor(listCobBO)
	updateCobProc := processor.NewUpdateCobProcessor(updateCobBO)
	deleteCobProc := processor.NewDeleteCobProcessor(deleteCobBO)
	createCobVProc := processor.NewCreateCobVProcessor(createCobVBO)
	getCobVProc := processor.NewGetCobVProcessor(getCobVBO)
	listCobVProc := processor.NewListCobVProcessor(listCobVBO)
	updateCobVProc := processor.NewUpdateCobVProcessor(updateCobVBO)
	simPayProc := processor.NewSimulatePaymentProcessor(simulatePayBO)
	createDevProc := processor.NewCreateDevolucaoProcessor(createDevBO)
	getDevProc := processor.NewGetDevolucaoProcessor(getDevBO)

	// Controllers
	createCobCtrl := controller.NewCreateCobController(createCobProc)
	getCobCtrl := controller.NewGetCobController(getCobProc)
	listCobCtrl := controller.NewListCobController(listCobProc)
	updateCobCtrl := controller.NewUpdateCobController(updateCobProc)
	deleteCobCtrl := controller.NewDeleteCobController(deleteCobProc)
	createCobVCtrl := controller.NewCreateCobVController(createCobVProc)
	getCobVCtrl := controller.NewGetCobVController(getCobVProc)
	listCobVCtrl := controller.NewListCobVController(listCobVProc)
	updateCobVCtrl := controller.NewUpdateCobVController(updateCobVProc)
	simPayCtrl := controller.NewSimulatePaymentController(simPayProc)
	createDevCtrl := controller.NewCreateDevolucaoController(createDevProc)
	getDevCtrl := controller.NewGetDevolucaoController(getDevProc)

	// Routes — cob (immediate charge)
	s.mux.HandleFunc("POST /cob", createCobCtrl.Handle)
	s.mux.HandleFunc("PUT /cob/{txid}", createCobCtrl.Handle)
	s.mux.HandleFunc("GET /cob", listCobCtrl.Handle)
	s.mux.HandleFunc("GET /cob/{txid}", getCobCtrl.Handle)
	s.mux.HandleFunc("PATCH /cob/{txid}", updateCobCtrl.Handle)
	s.mux.HandleFunc("DELETE /cob/{txid}", deleteCobCtrl.Handle)

	// Routes — cobv (charge with due date)
	s.mux.HandleFunc("PUT /cobv/{txid}", createCobVCtrl.Handle)
	s.mux.HandleFunc("GET /cobv", listCobVCtrl.Handle)
	s.mux.HandleFunc("GET /cobv/{txid}", getCobVCtrl.Handle)
	s.mux.HandleFunc("PATCH /cobv/{txid}", updateCobVCtrl.Handle)

	// Routes — simulate payment and refunds
	s.mux.HandleFunc("POST /cob/simulate", simPayCtrl.Handle)
	s.mux.HandleFunc("PUT /cob/{e2eid}/devolucao/{id}", createDevCtrl.Handle)
	s.mux.HandleFunc("GET /cob/{e2eid}/devolucao/{id}", getDevCtrl.Handle)

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	fmt.Printf("PIX PSP Simulator rodando em %s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}

// buildRepositories selects memory or SQL-backed implementations based on cfg.Storage.
func buildRepositories(cfg Config) (interfaces.CobRepository, interfaces.CobVRepository, interfaces.PixRepository) {
	if cfg.Storage == "sql" {
		db, err := sqlrepo.NewDB(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if err := sqlrepo.RunMigrations(db); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
		return sqlrepo.NewCobRepository(db),
			sqlrepo.NewCobVRepository(db),
			sqlrepo.NewPixRepository(db)
	}

	return memory.NewCobRepository(),
		memory.NewCobVRepository(),
		memory.NewPixRepository()
}
