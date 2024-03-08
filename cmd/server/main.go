package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/kameshsampath/go-cqrs-demo/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {

	//Setup Logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(setDBToContext)

	//Routes
	e.GET("/", routes.ListAll)
	e.GET("/:status", routes.ListByStatus)
	e.POST("/", routes.Create)
	e.PATCH("/:id", routes.Update)
	e.DELETE("/:id", routes.Delete)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// setDBToContext middleware initializes and sets the *gorm.DB to the echo context
func setDBToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Init/Setup DB and add it to the echo context
		var dbFile string
		cwd, _ := os.Getwd()
		if s, ok := os.LookupEnv("DB_FILE"); !ok {
			dbFile = path.Join(cwd, "data", "todo-test.db")
		} else {
			dbFile = s
		}
		log.Debugf("Using DB %s", dbFile)
		db, err := dao.New(dbFile)
		if err != nil {
			log.Fatal(err)
		}
		c.Set("DB", db)
		if err != nil {
			return err
		}
		return next(c)
	}
}
