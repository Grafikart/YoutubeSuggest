package main

import (
	"fmt"
	"github.com/Grafikart/YoutubeSuggest/api"
	"log"
	"os"
)

func main() {
	app, err := api.NewAPI(os.Getenv("GOOGLE_API_KEY"))
	catch(err)
	s, err := app.Subscriptions("UCj_iGliGCkLcHSZ8eqVNPDQ")
	catch(err)
	fmt.Printf("%d abonnements", len(s))
}

func catch(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
