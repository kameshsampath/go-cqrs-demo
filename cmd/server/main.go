package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/kameshsampath/go-cqrs-demo/routes"
	"github.com/kameshsampath/go-cqrs-demo/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var cfg *config.Config

func main() {

	cfg = config.New()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(setDBToContext)
	e.Use(setClientToContext)

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
		port := 8085
		if p, ok := os.LookupEnv("APP_PORT"); ok {
			port, _ = strconv.Atoi(p)
		}
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server,%s", err)
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
		cwd, _ := os.Getwd()
		dbFile := path.Join(cwd, cfg.DBFile)
		log.Debugf("Using DB %s", dbFile)
		db, err := dao.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		c.Set("DB", db)
		return next(c)
	}
}

// setClientToContext creates the Redpanda client and set it to the echo context
func setClientToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		client, err := utils.NewClient(cfg)
		if err != nil {
			return err
		}
		c.Set("RP_CLIENT", client)
		return next(c)
	}
}
