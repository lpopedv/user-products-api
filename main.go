package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/lpopedv/user-products-api/configs"
	"github.com/lpopedv/user-products-api/internal/entity"
	"github.com/lpopedv/user-products-api/internal/infra/database"
	"github.com/lpopedv/user-products-api/internal/infra/webserver/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lpopedv/user-products-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Products API
// @version         1.0
// @description     Product API with authentication

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in   header
// @name Authorization
func main() {
	cfg, err := configs.LoadConfig(".")

	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {


		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, cfg.TokenAuth, cfg.JWTExpiresIn)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.FindById)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/users", userHandler.Create)
	r.Post("/sessions", userHandler.Sessions)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
