package server

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"xmr-remote-nodes/frontend"
	"xmr-remote-nodes/handler"
	"xmr-remote-nodes/internal/config"
	"xmr-remote-nodes/internal/cron"
	"xmr-remote-nodes/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the WebUI",
	Long:  `This command will run HTTP server for APIs and WebUI.`,
	Run: func(_ *cobra.Command, _ []string) {
		serve()
	},
}

func serve() {
	appCfg := config.AppCfg()
	// connect to DB
	if err := database.ConnectDB(); err != nil {
		panic(err)
	}

	// run db migrations
	if err := database.MigrateDb(database.GetDB()); err != nil {
		panic(err)
	}

	// Define Fiber config & app.
	app := fiber.New(fiberConfig())

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

	handler.V1Api(app)
	app.Use("/", filesystem.New(filesystem.Config{
		Root: frontend.SvelteKitHandler(),
		// NotFoundFile: "index.html",
	}))

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start a cleanup cron-job
	if !fiber.IsChild() {
		cronRepo := cron.New()
		go cronRepo.RunCronProcess()
	}

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		slog.Info("Shutting down HTTP server...")
		_ = app.Shutdown()
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		slog.Error(fmt.Sprintf("Server is not running! error: %v", err))
	}
}

func fiberConfig() fiber.Config {
	return fiber.Config{
		Prefork:     config.AppCfg().Prefork,
		ProxyHeader: config.AppCfg().ProxyHeader,
		AppName:     "ditatompel's XMR Nodes HTTP server",
	}
}
