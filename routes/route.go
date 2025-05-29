package routes

import (
	"database/sql"
	"net/http"

	"project_minyak/middleware"
	"project_minyak/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRoutes(db *sql.DB, gormDB *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/signup", services.SignUp(gormDB)).Methods("POST")
	router.HandleFunc("/login", services.LoginHandler(db)).Methods("POST")
	router.HandleFunc("/midtrans/webhook", services.MidtransWebhookHandler(gormDB)).Methods("POST")

	// ADMIN ROUTES
	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Handle("/create-account", middleware.AdminMiddleware(http.HandlerFunc(services.CreateAccount(gormDB)))).Methods("POST")
	adminRoutes.Handle("/viewuser", middleware.AdminMiddleware(http.HandlerFunc(services.ViewUsers(gormDB)))).Methods("GET")
	adminRoutes.Handle("/delete-user", middleware.AdminMiddleware(http.HandlerFunc(services.DeleteUser(gormDB)))).Methods("DELETE")

	// SALES ROUTES
	salesRoutes := router.PathPrefix("/sales").Subrouter()
	salesRoutes.Use(middleware.SalesMiddleware)
	salesRoutes.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		services.InsertProductAndStock(w, r, gormDB)
	}).Methods("POST")
	router.HandleFunc("/stocks", services.GetStock(gormDB)).Methods("GET")
	salesRoutes.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		services.UpdateProductStock(w, r, gormDB)
	}).Methods("PUT")
	salesRoutes.HandleFunc("/rawmaterial", services.InsertRawMaterial(gormDB)).Methods("POST")
	salesRoutes.HandleFunc("/rawmaterial", services.GetRawMaterialsSorted(gormDB)).Methods("GET")

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

	customerRoutes.HandleFunc("/transactions/summary", services.ViewTransactionSummary(gormDB)).Methods("GET")
	customerRoutes.HandleFunc("/transactions/detail", services.ViewTransactionDetailByID(gormDB)).Methods("GET")
	customerRoutes.HandleFunc("/checkout", services.CheckoutHandler(gormDB)).Methods("POST")
	customerRoutes.HandleFunc("/cart", services.AddToCart(gormDB)).Methods("POST")         // Tambah item ke cart
	customerRoutes.HandleFunc("/cart", services.GetUserCart(gormDB)).Methods("GET")        // Ambil semua item cart user
	customerRoutes.HandleFunc("/cart", services.DeleteCartItems(gormDB)).Methods("DELETE") // Hapus item tertentu dari cart

	// ADMIN bisa akses semua route sales dan manager
	copyRoutes(adminRoutes, salesRoutes, middleware.AdminMiddleware)
	copyRoutes(adminRoutes, managerRoutes, middleware.AdminMiddleware)

	// MIDTRANS CALLBACK
	/*paymentRoutes := router.PathPrefix("/payment").Subrouter()
	paymentRoutes.HandleFunc("/midtrans/callback", func(w http.ResponseWriter, r *http.Request) {
		services.UpdateTransactionStatus(w, r, db)
	}).Methods("POST")*/

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
