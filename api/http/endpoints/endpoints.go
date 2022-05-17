package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	handlers "github.com/n0byk/loyalty/api/http/handlers/user"
	"github.com/n0byk/loyalty/api/http/middleware"
	"go.uber.org/zap"
)

var tokenAuth *jwtauth.JWTAuth

func InitEndpoints(logger *zap.Logger) chi.Router {

	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.LoggerMiddleware(logger))

		r.Post("/register", handlers.UserRegister)
		r.Post("/login", handlers.UserLogin)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Post("/orders", handlers.SetOrders)
			r.Get("/orders", handlers.GetOrders)

		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/balance", handlers.GetBalance)
			r.Post("/balance/withdraw", handlers.PostWithdraw)
			r.Get("/balance/withdrawals", handlers.GetWithdraws)

		})
	})

	// kubernetes Probe
	r.Route("/k8s", func(r chi.Router) {
		//https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
		r.Get("/livenessProbe", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		//https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
		r.Get("/readinessProbess", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
	})

	return r
}
