package main

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	backend_config "github.com/ponrove/ponrove-backend/pkg/config"
	frontend_config "github.com/ponrove/ponrove-frontend/pkg/config"
	"github.com/ponrove/ponrove-frontend/pkg/webclient"
	"github.com/ponrove/ponrunner"
	"github.com/rs/zerolog/log"
)

var bundles = []ponrunner.APIBundle{
	hub.Register,
	ingestion.Register,
}

func main() {
	var err error
	cfg := configura.Merge(frontend_config.New(), backend_config.New())

	// Add default logger to the context, which all http handlers derive their context (and logger) from.
	ctx := log.Logger.WithContext(context.Background())

	router := chi.NewRouter()

	// Start the runtime with the provided configuration and API bundles.
	err = ponrunner.Start(ctx, cfg, router, func(c configura.Config, r chi.Router, a huma.API) error {
		err := ponrunner.RegisterAPIBundles(c, a, bundles...)
		if err != nil {
			return err
		}

		err = webclient.Register(c, r)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start runtime")
	}
}
