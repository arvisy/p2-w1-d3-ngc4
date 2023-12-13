package main

import (
	"log"
	"net/http"
	"ngc4/config"
	"ngc4/handler"

	"github.com/julienschmidt/httprouter"
)

func main() {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	router := httprouter.New()

	router.GET("/avengers/inventory", handler.GetInventory)
	router.GET("/avengers/inventory/:id", handler.GetInventoryByID)
	router.POST("/avengers/inventory", handler.CreateInventory)
	router.DELETE("/avengers/inventory/:id", handler.DeleteInventoryByID)
	router.PUT("/avengers/inventory/:id", handler.UpdateInventoryID)

	router.GET("/avengers/crimeevent", handler.GetCrimeEvent)
	router.GET("/avengers/crimeevent/:id", handler.GetCrimeEventByID)
	router.POST("/avengers/crimeevent", handler.CreateCrimeEvent)
	router.DELETE("/avengers/crimeevent/:id", handler.DeleteCrimeEventByID)
	router.PUT("/avengers/crimeevent/:id", handler.UpdateCrimeEventByID)

	router.GET("/avengers/heroes", handler.GetHeroes)
	router.GET("/avengers/heroes/:id", handler.GetHeroesByID)
	router.POST("/avengers/heroes", handler.CreateHero)
	router.DELETE("/avengers/heroes/:id", handler.DeleteHeroByID)
	router.PUT("/avengers/heroes/:id", handler.UpdateHeroByID)

	router.GET("/avengers/villain", handler.GetVillain)
	router.GET("/avengers/villain/:id", handler.GetVillainByID)
	router.POST("/avengers/villain", handler.CreateVillain)
	router.DELETE("/avengers/villain/:id", handler.DeleteVillainByID)
	router.PUT("/avengers/villain/:id", handler.UpdateVillainByID)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
