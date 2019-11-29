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
	VehiclePlateNumber int    `required:"true" envconfig:"VEHICLE_PLATE_NUMBER"`
	NoAlertString      string `required:"true" envconfig:"NO_ALERT_STRING"`
}

type Notifier interface {
	notify() bool
}

func main() {
  var notifierParameter string

	err := godotenv.Load("/go/bin/.env")
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	var e Environment
	err = envconfig.Process("poundcheck", &e)

	if errParse != nil {
		log.Fatalf("envconfig.Process: %w", err.error)
	}

	app := cli.NewApp()
	app.Name = "paris-pound-check"
	app.Usage = "Check if you vehicle has been impounded"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "notifier, n",
			Value:       "none",
			Usage:       "Chose a notifier. Supported values : slack, sms",
			Destination: &notifierParameter,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "check if vehicle has been impounded",
			Action: func(c *cli.Context) error {
				var isImpounded bool

				notifer = getNotifier(notifierParameter)

				isImpounded = isVehicleImpounded(e)
				if isImpounded {
					log.Println("Vehicle was impounded, sending notification")

					notify(e, notifier)
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
				if !isNotifierValid(notifier) {
					log.Println("Please specify a notifier, invalid notifier specified")
					return nil
				}

				notify(notifier)
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func isVehicleImpounded(e Environment) bool {

	requestUrl := getPoundUrl(e)

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestUrl, nil)
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
    log.Fatal(err)
  }
  log.Println(fmt.Sprintf("Visiting %v", GetPoundUrl()))

	// Check for maintenance mode (happens a lot)
	if strings.Contains(string(body), "maintenance") {
		log.Println("Pound website in maintenance mode, checking later.")
		return false
	}

	if strings.Contains(string(body), e.NoAlertString) {
		return false
	}

	return true
}

func getPoundUrl(e Environment) string {
	return fmt.Sprintf("%v?immatriculation=%v&action_rechercher=", e.ParisPoundUrl, e.VehiclePlateNumber)
}


func getNotifier(  string) (Notifier,error) {

  switch notifierParameter {
  case "slack":
    return makeSlackNotifier()
  }

  return
}

func notify(e Environment, notifier string) {
	switch notifier {
	case "slack":
		if slackToken == "" || slackChannel == "" {
			log.Fatal("Missing environment variables for notifier Slack.")
		}
		message := GetNotificationMessage(getPoundUrl(e), vehiclePlateNumber)
		SendMessage(message)
	}
}
