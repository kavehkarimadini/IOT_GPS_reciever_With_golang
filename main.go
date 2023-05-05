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
	config := &serial.Config{Name: "COM4", Baud: 9600}
	ser, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	// set up Telegram bot
	bot, err := tg.NewBotAPI("5835596784:AAFxweJOUyXtIfhuvZVXnhghaJs-S47wYHg")
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
			_, err := ser.Write([]byte("$GPGGA\r\n"))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
			buf := make([]byte, 1024)
			n, err := ser.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(buf[:n]), "+++++++++++++++++++++++++++++++++++++++++++++++++")
			var foundStart bool
			var currentSentence string
			for _, r := range buf {
				if r == '$' {
					foundStart = true
					currentSentence = ""
				}

				if foundStart {
					currentSentence += string(r)
					if r == '\n' {
						if strings.HasPrefix(currentSentence, "$GPGGA") {
							fmt.Println(currentSentence)
							break
						}
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
