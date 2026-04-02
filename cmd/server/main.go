package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gewall/short-url/internal/handler"
	mw "github.com/gewall/short-url/internal/middleware"
	repository "github.com/gewall/short-url/internal/repository/postgres"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/internal/worker"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/oschwald/geoip2-golang/v2"
)

func main() {
	err := pkg.NewConfig().LoadConfig()
	if err != nil {
		panic(err)
	}

	geodb, err := geoip2.Open(os.Getenv("GEOIP_LOC"))
	if err != nil {
		panic(fmt.Errorf("failed to open geoip db: %w", err))
	}
	defer geodb.Close()

	db, err := pkg.NewPostgres()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OKE"))
	})

	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)

	refTokRepo := repository.NewRefreshTokenRepo(db)
	authSvc := service.NewAuthService(refTokRepo, userRepo)
	authHdl := handler.NewAuthHandler(authSvc)

	linkRepo := repository.NewLinkRepo(db)
	linkSvc := service.NewLinkService(linkRepo)
	linkHdl := handler.NewLinkHandler(linkSvc)

	clickRepo := repository.NewClickRepo(db)
	cw := worker.NewClickWorker(4, 100, service.RedirectProcessJob, clickRepo)
	redirectSvc := service.NewRedirectService(geodb, linkRepo, clickRepo, cw)
	redirectHdl := handler.NewRedirectHandler(redirectSvc)

	analyticSvc := service.NewAnalyticService(clickRepo)
	analyticHdl := handler.NewAnalyticHandler(analyticSvc)

	apiRouter := chi.NewRouter()
	apiRouter.Use(mw.JWTMiddleware)
	apiRouter.Route("/users", func(r chi.Router) {
		r.Get("/", userHdl.FindUsers)
		r.Post("/", userHdl.CreateUser)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", userHdl.FindUserByID)
			r.Delete("/", userHdl.DeleteUser)
		})
	})
	apiRouter.Route("/links", func(r chi.Router) {
		r.Get("/", linkHdl.FindAllByUserId)
		r.Post("/", linkHdl.CreateLink)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", linkHdl.FindById)
			r.Patch("/", linkHdl.UpdateLink)
			r.Delete("/", linkHdl.DeleteLink)
			r.Route("/analytics", func(r chi.Router) {

				r.Get("/", analyticHdl.Analytics)
				r.Get("/clicks", analyticHdl.TimeSeries)
				r.Get("/country", analyticHdl.Country)
			})
		})
	})

	authRouter := chi.NewRouter()
	authRouter.Route("/", func(r chi.Router) {
		r.Post("/signup", authHdl.SignUp)
		r.Post("/signin", authHdl.SignIn)
		r.Post("/refresh-token", authHdl.RefreshToken)
	})

	redirectRouter := chi.NewRouter()
	redirectRouter.Get("/{code}", redirectHdl.Redirect)

	r.Mount("/auth", authRouter)
	r.Mount("/api", apiRouter)
	r.Mount("/", redirectRouter)

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	http.ListenAndServe(":8080", r)
}
