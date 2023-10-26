package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type ServerImpl interface {
	Start()
}

type server struct {
	app *fiber.App
	cfg config.ConfigImpl
	db  *sqlx.DB
}

func NewServer(cfg config.ConfigImpl, db *sqlx.DB) ServerImpl {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(
			fiber.Config{
				AppName:      cfg.App().Name(),
				BodyLimit:    cfg.App().BodyLimit(),
				ReadTimeout:  cfg.App().ReadTimeout(),
				WriteTimeout: cfg.App().WriteTimeout(),
				JSONEncoder:  json.Marshal,
				JSONDecoder:  json.Unmarshal,
			},
		),
	}
}

func (s *server) Start() {

	// Middlewrare
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	// Modules
	v1 := s.app.Group("/v1")

	modules := InitModule(v1, s, middlewares)

	modules.MonitorModule()
	modules.UsersModule()
	modules.AppinfoModule()
	modules.FilesModule().Init()
	modules.ProductsModule().Init()
	modules.OrdersModule()

	s.app.Use(middlewares.RouterCheck())

	// Gacefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Shutting down server...")
		_ = s.app.Shutdown()
	}()

	// show URL APP
	log.Println("Server is running on port", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
