package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/n0byk/loyalty/api"
	handlers "github.com/n0byk/loyalty/api/http/handlers/user"
	"go.uber.org/zap"
)

var tokenAuth *jwtauth.JWTAuth

func InitEndpoints(logger *zap.Logger) chi.Router {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(api.LoggerMiddleware(logger))
		// r.Use(jwtauth.Verifier(tokenAuth))
		// r.Use(jwtauth.Authenticator)
		// r.Use(jwtauth.Verifier(jwtauth.New("HS256", []byte("secret"), nil)))
		// r.Use(jwtauth.Authenticator)
		r.Post("/api/user/register", handlers.UserRegister)
		r.Post("/api/user/login", handlers.UserLogin)
	})

	// kubernetes
	r.Route("/k8s", func(r chi.Router) {
		//https: //kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
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
