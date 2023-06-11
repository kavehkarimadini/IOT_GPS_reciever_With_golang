package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tarm/serial"
)

func main() {
	// set up serial port for GPS receiver
	config := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600}
	ser, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	// set up Telegram bot
	bot, err := tg.NewBotAPI("Your own bot api")
	if err != nil {
		log.Fatal(err)
	}

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.Text == "gps" {
			// retrieve GPS data
			// It sends the gpgga command to the gps receiver
			_, err := ser.Write([]byte("$GPGGA\r\n"))
			if err != nil {
				log.Fatal(err)
			}
			// waits 1 second for the response
			time.Sleep(time.Second)
			// makes a list of type byte with size of 1024 bytes
			buf := make([]byte, 1024)
			// reads the response from the gps receiver
			n, err := ser.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(buf[:n]), "+++++++++++++++++++++++++++++++++++++++++++++++++")
			var foundStart bool
			var currentSentence string
			for _, r := range buf {
				// iterates through buf, if $ is found then it rewrites the current sentence to empty string
				if r == '$' {
					foundStart = true
					currentSentence = ""
				}
				// if $ is found then it adds the current sentence with the elements of buf accordingly
				if foundStart {
					currentSentence += string(r)
					//when it encounters \n then it will check if it starts with gpgga, if it is found then we have our desired sentence
					if r == '\n' {
						if strings.HasPrefix(currentSentence, "$GPGGA") {
							fmt.Println(currentSentence)
							break
						}
						// if gpgga is not at the beginning of the sentence then mark found as false that will make it stop adding
						foundStart = false
					}
				}
			}

			msg, err := nmea.Parse(currentSentence)
			if err != nil {
				fmt.Println("Error parsing NMEA sentence:", err)
				return
			}

			// Extract the latitude and longitude values
			// checks if the msg is a gga format
			if gga, ok := msg.(nmea.GGA); ok {
				lat := gga.Latitude
				lon := gga.Longitude
				fmt.Println(lat)
				fmt.Println(lon)
				msgt := tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, lon))
				_, err = bot.Send(msgt)
				if err != nil {
					log.Fatal(err, "**********************************************************************************")
				}
				fmt.Printf("Latitude: %f, Longitude: %f\n", lat, lon)
			} else {
				fmt.Println("Invalid NMEA sentence type")
			}
		}
	}
}
