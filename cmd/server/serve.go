package server

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ditatompel/xmr-remote-nodes/frontend"
	"github.com/ditatompel/xmr-remote-nodes/internal/config"
	"github.com/ditatompel/xmr-remote-nodes/internal/cron"
	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/ditatompel/xmr-remote-nodes/internal/handler"
	"github.com/ditatompel/xmr-remote-nodes/internal/handler/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the WebUI and APIs",
	Long:  `This command will run HTTP server for APIs and WebUI.`,
	Run: func(_ *cobra.Command, _ []string) {
		serve()
	},
}

func serve() {
	appCfg := config.AppCfg()
	if err := database.ConnectDB(); err != nil {
		slog.Error(fmt.Sprintf("[DB] %s", err.Error()))
		os.Exit(1)
	}

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	stopCron := make(chan struct{})
	if !fiber.IsChild() {
		// run db migrations
		if err := database.MigrateDb(database.GetDB()); err != nil {
			slog.Error(fmt.Sprintf("[DB] %s", err.Error()))
			os.Exit(1)
		}

		// run cron process
		cronRepo := cron.New()
		go cronRepo.RunCronProcess(stopCron)
	}

	// Define Fiber config & app.
	app := fiber.New(fiber.Config{
		Prefork:     appCfg.Prefork,
		ProxyHeader: appCfg.ProxyHeader,
		AppName:     "XMR Nodes Aggregator",
	})

	// recover
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	// logger middleware
	if appCfg.LogLevel == "DEBUG" {
		app.Use(logger.New(logger.Config{
			Format:     "${time} DEBUG [HTTP] ${status} - ${latency} ${method} ${path} ${queryParams} ${ip} ${ua}\n",
			TimeFormat: "2006/01/02 15:04:05",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     appCfg.AllowOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	app.Use("/assets", views.EmbedAssets())

	handler.V1Api(app)
	app.Use("/", filesystem.New(filesystem.Config{
		Root: frontend.SvelteKitHandler(),
		// NotFoundFile: "index.html",
	}))

	// go routine to capture system calls
	go func() {
		<-sigCh
		close(stopCron) // stop cron goroutine
		slog.Info("Shutting down HTTP server...")
		_ = app.Shutdown()

		// give time for graceful shutdown
		time.Sleep(1 * time.Second)
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		slog.Error(fmt.Sprintf("[HTTP] Server is not running! error: %v", err))
	}
}
