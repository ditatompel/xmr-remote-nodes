package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ditatompel/xmr-nodes/frontend"
	"github.com/ditatompel/xmr-nodes/handler"
	"github.com/ditatompel/xmr-nodes/internal/config"
	"github.com/ditatompel/xmr-nodes/internal/database"
	"github.com/ditatompel/xmr-nodes/internal/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
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

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve() {
	appCfg := config.AppCfg()
	// connect to DB
	if err := database.ConnectDB(); err != nil {
		panic(err)
	}

	// Define Fiber config & app.
	app := fiber.New(fiberConfig())

	// recover
	app.Use(recover.New(recover.Config{EnableStackTrace: appCfg.Debug}))

	// logger middleware
	if appCfg.Debug {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} ${queryParams} ${ip} ${ua}\n",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     appCfg.AllowOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// cookie
	app.Use(encryptcookie.New(encryptcookie.Config{Key: appCfg.SecretKey}))

	handler.AppRoute(app)
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
		cronRepo := repo.NewCron(database.GetDB())
		go cronRepo.RunCronProcess()
	}

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		fmt.Println("Shutting down HTTP server...")
		_ = app.Shutdown()
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		fmt.Printf("Server is not running! error: %v", err)
	}
}

func fiberConfig() fiber.Config {
	return fiber.Config{
		Prefork:     config.AppCfg().Prefork,
		ProxyHeader: config.AppCfg().ProxyHeader,
		AppName:     "ditatompel's XMR Nodes HTTP server " + AppVer,
	}
}
