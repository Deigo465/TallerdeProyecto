package web

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	"github.com/open-wm/blockehr/pkg/infrastructure"
	"github.com/open-wm/blockehr/pkg/interfaces"
	repository "github.com/open-wm/blockehr/pkg/repositories"
	"github.com/open-wm/blockehr/pkg/web/handler"
)

type roleRouter struct {
	router     *mux.Router
	roles      []entities.Role
	middleware func(auth handler.AuthenticatedHandler, roles ...entities.Role) handler.Handler
}

func (r *roleRouter) reg(method string, route string, handler handler.AuthenticatedHandler) {
	r.router.HandleFunc(route, r.middleware(handler, r.roles...)).Methods(method)
}

func (r *roleRouter) PathPrefix(prefix string, handler http.Handler) {
	ah := func(user *entities.User, w http.ResponseWriter, r *http.Request) {
		log.Println("user is authenticated", user)
		handler.ServeHTTP(w, r)
	}
	r.router.PathPrefix(prefix).HandlerFunc(r.middleware(ah, r.roles...))
}

const VERSION = "0.2.0"

func StartWebServer(port string) {
	r := mux.NewRouter()

	defineRoutes(r, false)
	r.Use(LogMiddleware)

	log.Println("Starting server on port " + port)
	log.Println("Version: " + VERSION)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func defineRoutes(r *mux.Router, isTest bool) {
	logger := infrastructure.NewSLogger("blockehr")
	slog.SetDefault(logger)

	// this should probably go in a config file or options
	databaseName := "./db/dev.sqlite3"
	var db *sql.DB
	var blockchain interfaces.BlockchainClient
	handler.SetBasePath("./pkg/web/views/")
	if isTest {
		db = infrastructure.SetupTestDatabase()
		handler.SetBasePath("../../web/views/")
		blockchain = infrastructure.NewBlockchain("../../")
	} else {
		db = infrastructure.NewDatabase(logger, databaseName)
		log.Println("Database connected")
		blockchain = infrastructure.NewBlockchain("./")
		log.Println("Blockchain connected")

		// or if you want to use a mock blockchain
		// blockchain = mock.NewInMemoryBlockchain()
	}

	// repositories
	userRepo := repository.NewUserRepository(logger, db)
	appointmentRepo := repository.NewAppointmentRepository(logger, db)
	fileRepo := repository.NewFileRepository(logger, db)
	recordRepo := repository.NewRecordRepository(logger, db)

	sessRepo := repository.NewSessionRepository(db)
	profileRepo := repository.NewProfileRepository(logger, db)

	// use cases
	profileUC := usecase.NewProfileUsecase(profileRepo, appointmentRepo, userRepo)
	userUC := usecase.NewUserUsecase(userRepo)
	appointmentUC := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)
	fileUC := usecase.NewFileUsecase(fileRepo)
	recordUC := usecase.NewRecordUsecase(recordRepo, profileRepo, fileRepo, blockchain)
	authUC := usecase.NewAuthUsecase(userRepo, sessRepo, profileRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC)
	docHandler := handler.NewDoctorHandler(profileUC, userUC)
	patientsHandler := handler.NewPatientHandler(profileUC)
	appointmentHandler := handler.NewAppointmentHandler(appointmentUC)
	fileHandler := handler.NewFiletHandler(fileUC, recordUC)
	recordHandler := handler.NewRecordHandler(recordUC)

	// Routes
	r.HandleFunc("/", handler.Home).Methods(http.MethodGet)

	// Login
	r.HandleFunc("/login", authHandler.LoginForm).Methods(http.MethodGet)
	r.HandleFunc("/login", authHandler.LoginFormData).Methods(http.MethodPost)
	r.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost)

	sr := roleRouter{router: r, middleware: authHandler.WithUser, roles: []entities.Role{entities.STAFF}}

	// STAFF Views
	sr.reg(http.MethodGet, "/doctors", docHandler.List)
	sr.reg(http.MethodGet, "/patients", patientsHandler.List)

	// Staff and Doctor
	r.HandleFunc("/appointments", authHandler.WithUser(appointmentHandler.Calendar)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	enableCORS(api)

	SetAPIRoutes(api, authHandler, fileHandler, patientsHandler, docHandler, recordHandler, appointmentHandler)

	// 404
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		handler.ExecLayout(w, "404.html.tmpl", nil)
	})
	r.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) { handler.ExecLayout(w, "404", nil) })

	// serve static files in the public folder
	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("public/css/"))))
	r.PathPrefix("/img").Handler(http.StripPrefix("/img", http.FileServer(http.Dir("public/img/"))))
	r.PathPrefix("/js").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("public/js/"))))
	r.PathPrefix("/pdf").Handler(http.StripPrefix("/pdf", http.FileServer(http.Dir("public/pdf/"))))

	dr := roleRouter{router: r, middleware: authHandler.WithUser, roles: []entities.Role{entities.DOCTOR}}
	dr.PathPrefix("/downloads/", http.StripPrefix("/downloads", http.FileServer(http.Dir("storage/"))))
}
