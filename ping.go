/**
Created by Suko Widodo
**/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sparrc/go-ping"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	go tryping()
	fmt.Scanln()
}

func bottelegram(ip string) {
	bot, err := tgbotapi.NewBotAPI("XXXX.XXXXXX")
	if err != nil {
		log.Panic(err)
	}

	listuser := []int64{
		//-1001246507497, //id user chat group -> Send to group
		176274178, //id user chat -> Send to user
	}

	for i := 0; i < len(listuser); i++ {
		msg := tgbotapi.NewMessage(listuser[i], ip)
		bot.Send(msg)
	}

}

//Servers : Struct list of server
type Servers struct {
	Server []string `json:"server"`
}

//TryPing: func to run ping
func tryping() {
	/** How to get the current working directory in golang  **/
	/*dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)*/
	/** -------------------------------------------------**/

	// Open our jsonFile
	jsonFile, err := os.Open("pinger/server.json") // as a path of your json file //{"server": ["xxxx","xxxx"]}
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened server.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var servers Servers

	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &servers)

	timeout := flag.Duration("t", time.Second*2, "")
	interval := flag.Duration("i", time.Second*5, "")

	privileged := flag.Bool("privileged", true, "")
	for {
		for i := 0; i < len(servers.Server); i++ {
			pinger, err := ping.NewPinger(servers.Server[i])
			if err != nil {
				panic(err)
			}

			pinger.OnFinish = func(stats *ping.Statistics) {
				if stats.PacketsRecv > 0 {
					/** running on windows **/
					fmt.Fprintf(color.Output, "PC Server %s is UP \n", color.GreenString(servers.Server[i]))
					/** running on Linux **/
					//fmt.Printf(color.GreenString("%s Nyalah\n"), servers.Server[i])
				} else {
					/** running on windows **/
					fmt.Fprintf(color.Output, "PC Server %s is DOWN\n", color.RedString(servers.Server[i]))
					/** running on linux **/
					//fmt.Printf(color.RedString("%s Mati\n"), servers.Server[i])
					var kata = []string{"PC Server", servers.Server[i], "Mati"}
					var pesan = strings.Join(kata, " ")
					bottelegram(pesan) //Send to bot telegram
				}
			}

			pinger.SetPrivileged(*privileged)

			//fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
			pinger.Interval = *interval
			pinger.Timeout = *timeout
			pinger.Run()
		}

		time.Sleep(time.Second)
	}

}
