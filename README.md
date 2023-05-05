# IOT_GPS_reciever_With_golang
This repository is for iot purpose that is implemented using th GO language
This is a Go program that uses a GPS receiver connected to a serial port to retrieve location data and send the latitude and longitude information to a Telegram bot. Here is a breakdown of the code:

* The program imports several Go packages, including "fmt", "log", "strings", "time", "github.com/adrianmo/go-nmea", and "github.com/go-telegram-bot-api/telegram-bot-api".

* The program sets up a serial port configuration for the GPS receiver, with a name of "COM4" and a baud rate of 9600. Then, it opens the serial port using the "serial.OpenPort" function.

* The program sets up a Telegram bot using the bot's API token, and creates an update configuration with a timeout of 60 seconds. Then, it gets a channel of updates from the bot using the "bot.GetUpdatesChan" function.

* The program loops through each update received from the Telegram bot, ignoring any non-message updates.

* If the message text is "gps", the program sends a "$GPGGA" command to the GPS receiver over the serial port, reads the response from the GPS receiver, and parses the NMEA sentence to extract the latitude and longitude values.

* The program creates a new Telegram message with a URL to the location on Google Maps using the latitude and longitude values, and sends it to the chat ID of the message sender using the "bot.Send" function.

* Finally, the program prints the latitude and longitude values to the console.
