package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	slack_api "iamprashant.in/apps/slack-app/apis/slack-api"
	"iamprashant.in/apps/slack-app/configs"
	"iamprashant.in/apps/slack-app/internal/services"
	commons "iamprashant.in/apps/slack-app/pkg/commons"
	connectors "iamprashant.in/apps/slack-app/pkg/connectors"
)

type AppWrapper struct {
	E         *gin.Engine
	Cfg       *configs.AppConfig
	Postgres  connectors.PostgresConnector
	Logger    commons.Logger
	Closeable []func(context.Context) error
}

func main() {

	ctx := context.Background()
	appRunner := AppWrapper{E: gin.New()}

	// resolving configuration
	err := appRunner.ResolveConfig()
	if err != nil {
		panic(err)
	}
	// logging
	appRunner.Logging()

	// adding all connectors
	appRunner.AllConnectors()

	// init
	err = appRunner.Init(ctx)
	if err != nil {
		panic(err)
	}
	// add all routers
	appRunner.AllRouters()
	if err != nil {
		panic(err)
	}

	// running all migrator
	err = appRunner.E.Run()
	if err != nil {
		panic(err)
	}

	defer appRunner.Close(ctx)

}

// all middleware
func (g *AppWrapper) AllMiddlewares() {
	g.RecoveryMiddleware()
}

// Recovery middleware
func (g *AppWrapper) RecoveryMiddleware() {
	g.Logger.Info("Added Default Recovery middleware to the application.")
	g.E.Use(gin.Recovery())
}

func (g *AppWrapper) AllRouters() {
	g.SlackApiRoute()
}

func (g *AppWrapper) SlackApiRoute() {
	g.Logger.Info("Slack Api added to the application.")
	apiv1 := g.E.Group("/v1/")
	eventService := services.NewSlackEventService(g.Logger, g.Postgres)
	{
		apiv1.GET("/subscribe/", slack_api.NewSlackApi(g.Logger, eventService).Subscribe)
	}
	g.Logger.Info("Slack api and Services added to Server.")

}

func (app *AppWrapper) ResolveConfig() error {
	vConfig, err := configs.InitConfig()
	if err != nil {
		log.Fatalf("Unable to parse viper config to application configuration : %v", err)
		return err
	}

	cfg, err := configs.GetApplicationConfig(vConfig)
	if err != nil {
		log.Fatalf("Unable to parse viper config to application configuration : %v", err)
		return err
	}

	app.Cfg = cfg
	return nil
}

func (app *AppWrapper) Init(ctx context.Context) error {
	err := app.Postgres.Connect(ctx)
	if err != nil {
		app.Logger.Error("error while connecting to postgres.", err)
		return err
	}
	app.Closeable = append(app.Closeable, app.Postgres.Disconnect)
	return nil
}

func (app *AppWrapper) Logging() {
	aLogger := commons.NewApplicationLoggerWithOptions(
		&app.Cfg.Log,
		commons.Level(app.Cfg.LogLevel),
		commons.Name(app.Cfg.Name),
	)
	aLogger.InitLogger()
	aLogger.Info("Added Logger middleware to the application.")
	app.Logger = aLogger
}

func (g *AppWrapper) AllConnectors() {
	postgres := connectors.NewPostgresConnector(&g.Cfg.PostgresConfig, g.Logger)
	g.Postgres = postgres
}

// closer for app runner
func (app *AppWrapper) Close(ctx context.Context) {
	if len(app.Closeable) > 0 {
		app.Logger.Debug("there are closeable references to closed")
		for _, closeable := range app.Closeable {
			err := closeable(ctx)
			if err != nil {
				app.Logger.Errorf("error while closing %v", err)
			}
		}
	}
}
