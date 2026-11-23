package main

import (
	"fmt"
	"go-hexagonal-example/examples/core"
	"go-hexagonal-example/examples/db"
	"go-hexagonal-example/examples/utils"
	"log"
)

func main() {
	param := utils.ProvideFiberHttpServiceParams()
	fiberSrv := utils.InitializeHTTPService(param)

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal("Failed to start Database:", err)
	}
	db.DB = database

	repo := core.NewRepository(db.DB)
	srv := core.NewServices(repo)
	handler := core.NewHandlers(srv)

	fiberSrv.Get("/", handler.HandleGetTodoes)

	portString := fmt.Sprintf(":%v", param.Port)
	err = fiberSrv.Listen(portString)

	if err != nil {
		log.Fatal("Failed to start golang Fiber server:", err)
	}

}
