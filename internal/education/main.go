package main

import (
	"context"
	"net/http"

	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/maixuanbach174/online-course-app/internal/common/logs"
	"github.com/maixuanbach174/online-course-app/internal/common/server"
	"github.com/maixuanbach174/online-course-app/internal/education/ports"
	"github.com/maixuanbach174/online-course-app/internal/education/services"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application, err := services.NewApplication(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize application")
	}

	defer application.Close()

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(
				ports.NewHttpServer(application.App),
				router,
			)
		})
	}
}
