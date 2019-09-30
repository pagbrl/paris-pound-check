package main

import (
  "fmt"
  "log"
  "os"
  "strings"
  "io/ioutil"
  "net/http"

	"github.com/urfave/cli"
	"github.com/joho/godotenv"
)

var parisPoundUrl string
var vehiclePlateNumber string
var noAlertString string
var slackToken string
var slackChannel string

func init() {
	err := godotenv.Load("/go/bin/.env")
  if err != nil {
    log.Println("No .env file found, falling back to environment variables")
  }

	parisPoundUrl = os.Getenv("PARIS_POUND_URL")
  vehiclePlateNumber = os.Getenv("VEHICLE_PLATE_NUMBER")
  noAlertString = os.Getenv("NO_ALERT_STRING")
  slackToken = os.Getenv("SLACK_TOKEN")
  slackChannel = os.Getenv("SLACK_CHANNEL")

  if parisPoundUrl == "" || vehiclePlateNumber == "" || noAlertString == "" {
    log.Fatal("Required environment variables are missing.")
  }
}

func main() {
  var notifier string

  app := cli.NewApp()
  app.Name = "paris-pound-check"
  app.Usage = "Check if you vehicle has been impounded"

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "notifier, n",
      Value: "none",
      Usage: "Chose a notifier. Supported values : slack, sms",
      Destination: &notifier,
    },
  }

	app.Commands = []cli.Command{
    {
      Name:    "check",
      Aliases: []string{"c"},
      Usage:   "check if vehicle has been impounded",
      Action:  func(c *cli.Context) error {
        var isImpounded bool

        if (!IsNotifierValid(notifier)) {
          log.Println("Please specify a notifier, invalid notifier specified")
          return nil
        }

        isImpounded = IsVehicleImpounded()
        if (isImpounded) {
          log.Println("Vehicle was impounded, sending notification")

          Notify(notifier)
        } else {
          log.Println("Vehicle not impounded, nothing do to.")
        }

        return nil
      },
    },
    {
      Name: "test",
      Aliases: []string{"t"},
      Usage: "Test notifier settings",
      Action: func(c *cli.Context) error {
        log.Println("Sending test message")
        if (!IsNotifierValid(notifier)) {
          log.Println("Please specify a notifier, invalid notifier specified")
          return nil
        }

        Notify(notifier)
        return nil
      },
    },
  }

	err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
	}
}


func IsVehicleImpounded() bool {

  requestUrl := GetPoundUrl()

  client := &http.Client {}

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

  // Check for maintenance mode (happens a lot)
  if (strings.Contains(string(body), "maintenance")) {
    log.Println("Pound website in maintenance mode, checking later.")
    return false
  }

  if (strings.Contains(string(body), noAlertString)) {
    return false
  } else {
    return true
  }

}

func GetPoundUrl() string {
  return fmt.Sprintf("%v?immatriculation=%v&action_rechercher=", parisPoundUrl, vehiclePlateNumber)
}

func IsNotifierValid(notifier string) bool {
  if (notifier == "") {
    return false
  }

  switch notifier {
    case
        "sms",
        "slack":
        return true
    }
  return false
}

func Notify(notifier string) {
  switch notifier {
  case "slack":
    if slackToken == "" || slackChannel == "" {
      log.Fatal("Missing environment variables for notifier Slack.")
    }
    message := GetNotificationMessage(GetPoundUrl(), vehiclePlateNumber)
    SendMessage(message)
  }
}
