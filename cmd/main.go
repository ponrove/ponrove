package main

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/ponrove/ponrove-backend/pkg/api/hub"
	"github.com/ponrove/ponrove-backend/pkg/api/ingestion"
	backend_config "github.com/ponrove/ponrove-backend/pkg/config"
	"github.com/ponrove/ponrunner"
	"github.com/rs/zerolog/log"
)

var bundles = []ponrunner.APIBundle{
	hub.Register,
	ingestion.Register,
}

func main() {
	var err error
	backendCfg := backend_config.New()

	// Add default logger to the context, which all http handlers derive their context (and logger) from.
	ctx := log.Logger.WithContext(context.Background())

	router := chi.NewRouter()

	// Start the runtime with the provided configuration and API bundles.
	err = ponrunner.Start(ctx, backendCfg, router, bundles...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start runtime")
	}
}
