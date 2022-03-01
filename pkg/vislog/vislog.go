package vislog

// I need to find some more complicated examples.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// stuck on design decision...
// thinking..............
type parampair struct {
	Key   string
	Value string
}

type Vislogger struct {
	bulbid  string
	base    colour // base? needs a better name, default is a reserve word I think...
	colours map[float64]colour
}

type colour struct {
	name       string //do I need a name? yeah..
	hue        float64
	saturation float64
	brightness float64
	kelvin     int
	priority   int //will keep it simple "5" top priority "1" lowest
}

type status struct {
	Power      string  `json:"power"`
	Connected  bool    `json:"connected"`
	Brightness float64 `json:"brightness"`
	Color      struct {
		Hue        float64 `json:"hue"`
		Saturation float64 `json:"saturation"`
		Kelvin     int     `json:"kelvin"`
	} `json:"color"`
}

func NewBulb(bulbid string) *Vislogger {
	// should  CRUD for adding colours/default
	colours := map[float64]colour{
		0.0: {
			name:       "error",
			hue:        0.0,
			saturation: 1.0,
			brightness: 1.0,
			kelvin:     6000,
			priority:   5,
		},
		36.0: {
			name:       "warning",
			hue:        36.0,
			saturation: 1.0,
			brightness: 1.0,
			kelvin:     6000,
			priority:   4,
		},
		120.0: {
			name:       "success",
			hue:        120.0,
			saturation: 1.0,
			brightness: 1.0,
			kelvin:     6000,
			priority:   3,
		},
	}
	//base should have some sort of CRUD?
	base := colour{
		name:       "default",
		hue:        0.0,
		saturation: 0.0,
		brightness: 1.0,
		kelvin:     6000,
		priority:   1,
	}
	return &Vislogger{bulbid, base, colours}
}

func lnquery(bulbid string, capability string, http_method string, pairs []parampair) []byte {
	client := &http.Client{}
	parm := url.Values{}
	for _, pair := range pairs {
		parm.Add(pair.Key, pair.Value)
	}
	var lurl = fmt.Sprintf("https://api.lifx.com/v1/lights/id:%s/%s", bulbid, capability)
	req, err := http.NewRequest(http_method, lurl, nil)
	if err != nil {
		log.Println(err.Error())
	}
	req.URL.RawQuery = parm.Encode()
	var bearer = "Bearer " + os.Getenv("LIFXTOKEN")
	req.Header.Add("Authorization", bearer)
	resp, cerr := client.Do(req)
	if cerr != nil {
		log.Println(cerr)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func (v *Vislogger) priorityCheck(newPriority int) string {
	persist := "true"
	data := v.Status()
	log.Printf("Old Priority: %d\n", v.colours[data.Color.Hue].priority)
	if (data.Color.Saturation > 0.0) && (v.colours[data.Color.Hue].priority > newPriority) {
		log.Println("hello")
		persist = "false"
	}
	log.Printf("Persist is: %s\n", persist)
	return persist
}

func (v *Vislogger) Error() {
	persist := v.priorityCheck(v.colours[0.0].priority)
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("hue:%f saturation:1.0", v.colours[0.0].hue)}, {Key: "period", Value: "2.0"}, {Key: "cycles", Value: "5.0"}, {Key: "persist", Value: persist}}
	lnquery(v.bulbid, "effects/pulse", "POST", pairs)
}

func (v *Vislogger) Warning() {
	persist := v.priorityCheck(v.colours[36.0].priority)
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("hue:%f saturation:1.0", v.colours[36.0].hue)}, {Key: "period", Value: "2.0"}, {Key: "cycles", Value: "5.0"}, {Key: "persist", Value: persist}}
	lnquery(v.bulbid, "effects/pulse", "POST", pairs)
}

func (v *Vislogger) Success() {
	persist := v.priorityCheck(v.colours[120.0].priority)
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("hue:%f saturation:1.0", v.colours[120.0].hue)}, {Key: "period", Value: "2.0"}, {Key: "cycles", Value: "5.0"}, {Key: "persist", Value: persist}}
	lnquery(v.bulbid, "effects/pulse", "POST", pairs)
}

func (v *Vislogger) Normal() {
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("kelvin:%d", v.base.kelvin)}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *Vislogger) RGB(rgb string) {
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("rgb:%s saturation:1.0", rgb)}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *Vislogger) HEX(hex string) {
	if !strings.HasPrefix(hex, "#") {
		hex = fmt.Sprintf("#%s", hex)
	}
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("%s saturation:1.0", hex)}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *Vislogger) Party() {
	min := 0
	max := 255
	for {
		rgb := fmt.Sprintf("%d,%d,%d", rand.Intn(max-min)+min, rand.Intn(max-min)+min, rand.Intn(max-min)+min)
		pairs := []parampair{{Key: "color", Value: fmt.Sprintf("rgb:%s saturation:1.0", rgb)}}
		lnquery(v.bulbid, "state", "PUT", pairs)
		time.Sleep(500 * time.Millisecond)
	}
}

func (v *Vislogger) Dim() {
	pairs := []parampair{{Key: "color", Value: "brightness:0.5"}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *Vislogger) IsOn() bool {
	if v.Status().Power == "on" {
		return true
	} else {
		return false
	}
}

func (v *Vislogger) IsReady() bool {
	data := v.Status()
	return data.Connected
}

func (v *Vislogger) ShowAll() {
	body := lnquery(v.bulbid, "", "GET", []parampair{})
	fmt.Println(string(body))
}

func (v *Vislogger) Status() status {
	body := lnquery(v.bulbid, "", "GET", []parampair{})
	var data []status
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data[0]
}

//this delay thing is the most probable thing...api doesn't give back the
//proper result immediatly... so I need some sort of mechanism that locks
//changes until bulb is ready to receive.
// might be still a bit out of my experience level.

//func (v *Vislogger) doubleCheck(pairs []parampair) {
//	log.Println("Double Checking if the bulb is stuck")
//}
