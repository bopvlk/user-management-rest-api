package main

import (
	"log"

	"git.foxminded.com.ua/3_REST_API/interal/config"
	"git.foxminded.com.ua/3_REST_API/interal/infrastructure/datastore"
	"git.foxminded.com.ua/3_REST_API/interal/infrastructure/router"
	"git.foxminded.com.ua/3_REST_API/interal/registry"
	"github.com/labstack/echo/v4"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := datastore.NewDB(&config)

	r := registry.NewRegistry(db, &config)

	e := echo.New()
	e = router.NewRouter(e, r.NewAppController())

	log.Println("Server listen at http://localhost" + ":" + config.Addr)
	log.Fatalln(e.Start(":" + config.Addr))
}
