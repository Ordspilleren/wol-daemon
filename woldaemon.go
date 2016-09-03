package main

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/sabhiram/go-wol"
)

// BroadcastIP sets the IP that magic packets will be broadcast to
const BroadcastIP = "255.255.255.255"
// UDPPort sets the port on which packets will be sent to
const UDPPort = "9"
// BroadcastInterface sets the physical interface where packets will be sent from
const BroadcastInterface = ""

func main() {
	http.HandleFunc("/wake/", WakeRequest)

	http.ListenAndServe(":8033", nil)
}

// Machine defines the attributes of a PC 
type Machine struct {
	Name string
	IPAddr string
	MACAddr string
}

// WakeRequest handles HTTP requests to wake a machine
func WakeRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	URLPath := strings.Split(r.URL.Path, "/")
	pc := URLPath[2]
	machines := make(map[string]*Machine)
	machines["lasse-pc"] = &Machine{Name: "Lasse-PC", IPAddr: "192.168.1.2", MACAddr: "10:c3:7b:6b:f8:f2"}

	if val, ok := machines[pc]; ok {
		fmt.Println(val.IPAddr, val.MACAddr)

		wol.SendMagicPacket(val.MACAddr, BroadcastIP+":"+UDPPort, BroadcastInterface)
		fmt.Fprint(w, "Magic Packet sent to "+val.Name+"@"+val.MACAddr)
	}
}