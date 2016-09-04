package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"github.com/sabhiram/go-wol"
	"github.com/spf13/viper"
)

// BroadcastIP sets the IP that magic packets will be broadcast to
var BroadcastIP string
// UDPPort sets the port on which packets will be sent to
var UDPPort string
// BroadcastInterface sets the physical interface where packets will be sent from
var BroadcastInterface string

// Machine defines the attributes of a PC 
type Machine struct {
	Name string `mapstructure:"name"`
	IPAddr string `mapstructure:"ip"`
	MACAddr string `mapstructure:"mac"`
}

func main() {
	viper.SetConfigName("wol-daemon")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Config file could not be read: %s \n", err)
	}
	viper.SetDefault("BroadcastIP", "255.255.255.255")
	viper.SetDefault("UDPPort", "9")
	viper.SetDefault("BroadcastInterface", "")
	BroadcastIP = viper.GetString("BroadcastIP")
	UDPPort = viper.GetString("UDPPort")
	BroadcastInterface = viper.GetString("BroadcastInterface")

	http.HandleFunc("/wake/", WakeRequest)

	http.ListenAndServe(":8033", nil)
}

// WakeRequest handles HTTP requests to wake a machine
func WakeRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	URLPath := strings.Split(r.URL.Path, "/")
	pc := "machines."+URLPath[2]

	if viper.IsSet(pc) {
		var M Machine
		viper.UnmarshalKey(pc, &M)

		fmt.Println(M.Name, M.IPAddr, M.MACAddr)

		err := wol.SendMagicPacket(M.MACAddr, BroadcastIP+":"+UDPPort, BroadcastInterface)
		if err != nil {
			fmt.Fprint(w, "Could not send magic packet to "+M.MACAddr, err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			fmt.Fprint(w, "Magic Packet sent to "+M.Name+"@"+M.MACAddr)
			w.WriteHeader(http.StatusOK)
		}
	}
}