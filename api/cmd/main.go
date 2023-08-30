package main

import (
	"api/db"
	"api/internal/segment"
	"api/router"
	"context"
	"log"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not init db connection: %s", err)
	}
	segmentRepository := segment.NewRepository(dbConn.GetDB())
	segmentService := segment.NewService(segmentRepository)
	segmentHandler := segment.NewHandler(segmentService)

	segmentService.TtlService(context.Background())

	router.InitRouter(segmentHandler)
	err = router.Start("0.0.0.0:8081")
	if err != nil {
		log.Fatalf("%s", err)
		return
	}
}
