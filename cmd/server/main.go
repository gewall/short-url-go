package main

import (
	"net/http"

	"github.com/gewall/short-url/internal/handler"
	mw "github.com/gewall/short-url/internal/middleware"
	repository "github.com/gewall/short-url/internal/repository/postgres"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	err := pkg.NewConfig().LoadConfig()
	if err != nil {
		panic(err)
	}

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

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)

	refTokRepo := repository.NewRefreshTokenRepo(db)
	authSvc := service.NewAuthService(refTokRepo, userRepo)
	authHdl := handler.NewAuthHandler(authSvc)

	linkRepo := repository.NewLinkRepo(db)
	linkSvc := service.NewLinkService(linkRepo)
	linkHdl := handler.NewLinkHandler(linkSvc)

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
		})
	})

	authRouter := chi.NewRouter()
	authRouter.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHdl.SignUp)
		r.Post("/signin", authHdl.SignIn)
		r.Post("/refresh-token", authHdl.RefreshToken)
	})

	r.Mount("/", authRouter)
	r.Mount("/api", apiRouter)

	http.ListenAndServe(":8080", r)
}
