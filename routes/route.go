package routes

import (
	"database/sql"
	"net/http"

	"project_minyak/middleware"
	"project_minyak/services"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/signup", services.SignUp(db)).Methods("POST")
	router.HandleFunc("/login", services.LoginHandler(db)).Methods("POST")

	// ADMIN ROUTES
	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Handle("/create-account", middleware.AdminMiddleware(http.HandlerFunc(services.CreateAccount(db)))).Methods("POST")

	// SALES ROUTES
	salesRoutes := router.PathPrefix("/sales").Subrouter()
	salesRoutes.Use(middleware.SalesMiddleware) // if you have one
	salesRoutes.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		services.InsertProductAndStock(w, r, db)
	}).Methods("POST")
	salesRoutes.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		services.UpdateProductStock(w, r, db)
	}).Methods("PUT")
	salesRoutes.HandleFunc("/rawmaterial", func(w http.ResponseWriter, r *http.Request) {
		services.InsertRawMaterial(w, r, db)
	}).Methods("POST")
	salesRoutes.HandleFunc("/rawmaterial", func(w http.ResponseWriter, r *http.Request) {
		services.GetRawMaterialsSorted(w, r, db)
	}).Methods("GET")

	// MANAGER ROUTES
	managerRoutes := router.PathPrefix("/manager").Subrouter()
	managerRoutes.Use(middleware.ManagerMiddleware) // if you have one
	managerRoutes.HandleFunc("/gemini/analyze", services.AnalyzeAllProductsHandler(db)).Methods("POST")
	managerRoutes.HandleFunc("/transaction-view", func(w http.ResponseWriter, r *http.Request) {
		services.SalesRecap(w, r, db)
	}).Methods("GET")

	// CUSTOMER ROUTES
	customerRoutes := router.PathPrefix("/customer").Subrouter()
	customerRoutes.Use(middleware.CustomerMiddleware)

	customerRoutes.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		services.ViewTransaction(w, r, db)
	}).Methods("GET")

	// ADMIN bisa akses semua route sales dan manager
	copyRoutes(adminRoutes, salesRoutes, middleware.AdminMiddleware)
	copyRoutes(adminRoutes, managerRoutes, middleware.AdminMiddleware)

	// MIDTRANS CALLBACK
	paymentRoutes := router.PathPrefix("/payment").Subrouter()
	paymentRoutes.HandleFunc("/midtrans/callback", func(w http.ResponseWriter, r *http.Request) {
		services.UpdateTransactionStatus(w, r, db)
	}).Methods("POST")

	return router
}

func copyRoutes(targetRouter, sourceRouter *mux.Router, middleware func(http.Handler) http.Handler) {
	sourceRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		handler := route.GetHandler()
		if handler == nil {
			return nil
		}

		targetRouter.Handle(path, middleware(handler)).Methods(methods...)

		return nil
	})
}
