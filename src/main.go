package main

import (
	router2 "ad_service/http/router"
	"ad_service/infrastructure"
	"ad_service/interactor"
)

func main() {

	cassandraClient, err := infrastructure.NewCassandraSession()
	if err != nil {
		return
	}
	interactor2 := interactor.NewInteractor(cassandraClient)
	appHandler := interactor2.NewAppHandler()

	router := router2.NewRouter(appHandler)

	err = router.Run(":8093")
	if err != nil {
		return
	}

}
