package main

import (
	"book-shop/config"
	"book-shop/internal/app/http/handlers"
	"book-shop/internal/app/logger"
	"book-shop/internal/app/repositories"
	"book-shop/internal/app/services"
	"book-shop/internal/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	background := context.Background()
	loadConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("run: error load config %w", err)
	}

	err = logger.SetupLogger()
	if err != nil {
		return fmt.Errorf("run: error setup logger %w", err)
	}
	defer logger.Logger.Sync()

	err = RunMigrations(loadConfig.DSN, loadConfig.MigrationsPath)
	if err != nil {
		return fmt.Errorf("run: error run migration %w", err)
	}

	dbCon, err := postgres.Dial(background, loadConfig.DSN)
	if err != nil {
		return fmt.Errorf("run: error get connection %w", err)
	}

	bookRepository := repositories.NewBookRepository(dbCon)
	bookService := services.NewBookService(bookRepository)

	categoryRepository := repositories.NewCategoryRepository(dbCon)
	categoryService := services.NewCategoryService(categoryRepository)

	userRepository := repositories.NewUserRepository(dbCon)
	userService := services.NewUserService(userRepository)

	tokenRepository := repositories.NewTokenRepository(dbCon)
	jwtService := services.NewJWTService(tokenRepository)

	httpServer := handlers.NewHttpServer(bookService, categoryService, userService, jwtService)

	router := mux.NewRouter()

	//Сделал все по книгам, пока без проверки админ или не админ
	router.HandleFunc("/book/{book_id}", httpServer.GetBook).Methods(http.MethodGet)                              //получение книги
	router.HandleFunc("/books", httpServer.GetBooksByCategories).Methods(http.MethodGet)                          //получение страницы
	router.HandleFunc("/book", httpServer.CheckAdmin(httpServer.CreateBook)).Methods(http.MethodPost)             //Admin validation
	router.HandleFunc("/book/{book_id}", httpServer.CheckAdmin(httpServer.UpdateBook)).Methods(http.MethodPut)    //admin validation
	router.HandleFunc("/book/{book_id}", httpServer.CheckAdmin(httpServer.DeleteBook)).Methods(http.MethodDelete) //admin validation

	//Сделал все по категориям, пока без проверки админ или нет
	router.HandleFunc("/category", httpServer.CheckAdmin(httpServer.CreateCategory)).Methods(http.MethodPost)
	router.HandleFunc("/categories", httpServer.GetCategories).Methods(http.MethodGet)
	router.HandleFunc("/category/{category_id}", httpServer.CheckAdmin(httpServer.GetCategory)).Methods(http.MethodGet)
	router.HandleFunc("/category/{category_id}", httpServer.CheckAdmin(httpServer.UpdateCategory)).Methods(http.MethodPut)
	router.HandleFunc("/category/{category_id}", httpServer.CheckAdmin(httpServer.DeleteCategory)).Methods(http.MethodDelete)

	//Регистрация и аутентификация
	router.HandleFunc("/signup", httpServer.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", httpServer.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/logout", httpServer.CheckAuthorizedUser(httpServer.Logout)).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	jwtService.StartTokenCleanupScheduler()

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	//log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")

	return nil
}

func RunMigrations(dsn string, migrationsPath string) error {
	if dsn == "" {
		return errors.New("dsn is required")
	}

	if migrationsPath == "" {
		return errors.New("migrationsPath is required")
	}

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("problems with connection: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("problems migration execution: %w", err)
	}

	return nil
}
