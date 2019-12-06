package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli"
)

type Environment struct {
	ParisPoundUrl      string `required:"true" envconfig:"PARIS_POUND_URL"`
	VehiclePlateNumber string `required:"true" envconfig:"VEHICLE_PLATE_NUMBER"`
	NoAlertString      string `required:"true" envconfig:"NO_ALERT_STRING"`
}

type Notifier interface {
	Notify(detailsUrl string, vehiclePlateNumber string) bool
}

func main() {
  	var notifierParameter string

	err := godotenv.Load("/go/bin/.env")
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	var e Environment
	err = envconfig.Process("poundcheck", &e)
	if err != nil {
		log.Fatalf("envconfig.Process: %w", err)
	}

	app := cli.NewApp()
	app.Name = "paris-pound-check"
	app.Usage = "Check if you vehicle has been impounded"

	app.Flags = []cli.Flag {
		&cli.StringFlag{
		  Name:        "notifier, n",
		  Value:       "slack",
		  Usage:       "Chose a notifier. Supported values : slack",
		  Destination: &notifierParameter,
		},
	  }

	app.Commands = []*cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "check if vehicle has been impounded",
			Action: func(c *cli.Context) error {
				var isImpounded bool

				notifier := getNotifier(notifierParameter)

				poundUrl := getPoundUrl(e)

				isImpounded = isVehicleImpounded(poundUrl, e.NoAlertString)
				if isImpounded {
					log.Println("Vehicle was impounded, sending notification")

					notifier.Notify(poundUrl, e.VehiclePlateNumber)
				}
				log.Println("Vehicle not impounded, nothing do to.")

				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Test notifier settings",
			Action: func(c *cli.Context) error {
				log.Println("Sending test message")

				notifier := getNotifier(notifierParameter)
				notifier.Notify("test", "test")
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func isVehicleImpounded(poundUrl string, noAlertString string) bool {
	client := &http.Client{}

	req, err := http.NewRequest("GET", poundUrl, nil)
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
    	log.Fatal(err)
  	}
  	log.Println(fmt.Sprintf("Visiting %v", poundUrl))

	// Check for maintenance mode (happens a lot)
	if strings.Contains(string(body), "maintenance") {
		log.Println("Pound website in maintenance mode, checking later.")
		return false
	}

	if strings.Contains(string(body), noAlertString) {
		return false
	}

	return true
}

func getPoundUrl(e Environment) string {
	return fmt.Sprintf("%v?immatriculation=%v&action_rechercher=", e.ParisPoundUrl, e.VehiclePlateNumber)
}


func getNotifier(notifierParameter string) (Notifier) {

  switch notifierParameter {
  case "slack":
    return makeSlackNotifier()
  }

  return makeSlackNotifier()
}
