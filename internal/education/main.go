package main

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/logs"
	"github.com/maixuanbach174/online-course-app/internal/education/services"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application, err := services.NewApplication(ctx)
	if err != nil {
		return
	}

	defer application.Close()
}
