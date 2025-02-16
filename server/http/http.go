package http

import (
	"net/http"

	elect "github.com/aremxyplug-be/lib/bills/electricity"
	"github.com/aremxyplug-be/lib/bills/tvsub"
	"github.com/aremxyplug-be/lib/emailclient"
	otpgen "github.com/aremxyplug-be/lib/otp_gen"
	"github.com/aremxyplug-be/lib/telcom/airtime"
	"github.com/aremxyplug-be/lib/telcom/data"
	"github.com/aremxyplug-be/lib/telcom/edu"
	"github.com/aremxyplug-be/server/http/handlers"

	"github.com/aremxyplug-be/config"
	"github.com/aremxyplug-be/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Logger      *zap.Logger
	Store       db.DataStore
	Secrets     *config.Secrets
	EmailClient emailclient.EmailClient
	DataClient  *data.DataConn
	EduClient   *edu.EduConn
	Vtu         *airtime.AirtimeConn
	TvSub       *tvsub.TvConn
	ElectSub    *elect.ElectricConn
	Otp         *otpgen.OTPConn
}

func MountServer(config ServerConfig) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		Debug:            true,
	}).Handler)
	router.Use(setJSONContentType)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	// Get handlers
	httpHandler := handlers.NewHttpHandler(&handlers.HandlerOptions{
		Logger:      config.Logger,
		Store:       config.Store,
		Secrets:     config.Secrets,
		EmailClient: config.EmailClient,
		Data:        config.DataClient,
		Edu:         config.EduClient,
		VTU:         config.Vtu,
		TvSub:       config.TvSub,
		ElectSub:    config.ElectSub,
		Otp:         config.Otp,
	})

	// Routes
	// Health check
	router.Get("/health", healthCheck)

	router.Route("/api/v1", func(router chi.Router) {
		// SignUp
		router.Post("/signup", httpHandler.SignUp)
		// Login
		router.Post("/login", httpHandler.Login)
		// forgot password
		router.Post("/forgot-password", httpHandler.ForgotPassword)
		// reset password
		router.Patch("/reset-password", httpHandler.ResetPassword)

		router.Post("/send-otp", httpHandler.SendOTP)

		router.Post("/verify-otp", httpHandler.VerifyOTP)

		// test
		router.Post("/test", httpHandler.Testtoken)

		// Data Routes
		dataRoutes(router, httpHandler)
		// smile data routes
		smileDataRoutes(router, httpHandler)
		// spectranet data routes
		spectranetDataRoutes(router, httpHandler)

		// Edu Routes
		eduRoutes(router, httpHandler)

		//  Airtime Routes
		airtimeRoutes(router, httpHandler)

		// TvSubscription Routes
		tvSubscriptionRoutes(router, httpHandler)

		// Electricity bills routes
		electricityBillRoutes(router, httpHandler)

	})

	return router
}

func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.Data(w, r, []byte("Ok"))
}

func dataRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/data", func(router chi.Router) {
		router.Post("/", httpHandler.Data)
		router.Get("/", httpHandler.Data)
		router.Get("/{id}", httpHandler.GetDataInfo)
		router.Get("/transactions", httpHandler.GetDataTransactions)
	})
}

func smileDataRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/data/smile", func(router chi.Router) {
		router.Post("/", httpHandler.SmileData)
		router.Get("/", httpHandler.SmileData)
		router.Get("/{id}", httpHandler.GetSmileDataDetails)
		router.Get("/transactions", httpHandler.GetSmileTransactions)
	})
}

func spectranetDataRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/data/spectranet", func(router chi.Router) {
		router.Post("/", httpHandler.SpectranetData)
		router.Get("/", httpHandler.SpectranetData)
		router.Get("/{id}", httpHandler.GetSpecDataDetails)
		router.Get("/transactions", httpHandler.GetSpectranetTransactions)
	})
}

func eduRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/edu", func(router chi.Router) {
		router.Post("/", httpHandler.EduPins)
		router.Get("/", httpHandler.EduPins)
		router.Get("/{id}", httpHandler.GetDataInfo)
		router.Get("/transactions", httpHandler.GetEduTransactions)
	})
}

func airtimeRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/airtime", func(router chi.Router) {
		router.Post("/", httpHandler.Airtime)
		router.Get("/", httpHandler.Airtime)
		router.Get("/{id}", httpHandler.GetAirtimeInfo)
		router.Get("/transactions", httpHandler.GetAirtimeTransactions)
	})
}

func tvSubscriptionRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/tvsub", func(router chi.Router) {
		router.Post("/", httpHandler.TVSubscriptions)
		router.Get("/", httpHandler.TVSubscriptions)
		router.Get("/{id}", httpHandler.GetTvSubDetails)
		router.Get("/transactions", httpHandler.GetTvSubscriptions)
	})
}

func electricityBillRoutes(r chi.Router, httpHandler *handlers.HttpHandler) {
	r.Route("/electric-bill", func(router chi.Router) {
		router.Post("/", httpHandler.ElectricBill)
		router.Get("/", httpHandler.ElectricBill)
		router.Get("/{id}", httpHandler.GetElectricBillDetails)
		router.Get("/transactions", httpHandler.GetElectricBills)
	})
}
