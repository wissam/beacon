package vislog

// this is functioning as intended!
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

	"github.com/hofstadter-io/hof/lib/dotpath"
)

type parampair struct {
	Key   string
	Value string
}

type vislogger struct {
	bulbid  string
	base    colour // base? needs a better name, default is a reserve word I think...
	colours map[float64]colour
}

type colour struct {
	name       string //do I need a name? yeah..
	hue        float64
	saturation float64
	brightness float64
	kelvin     float64
	priority   int //will keep it simple "5" top priority "1" lowest
}

// let's see if this works...
func NewBulb(bulbid string) *vislogger {
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
	return &vislogger{bulbid, base, colours}
}

func lnquery(bulbid string, capability string, http_method string, pairs []parampair) []byte {
	log.Printf("%s", http_method)
	client := &http.Client{}
	parm := url.Values{}
	for _, pair := range pairs {
		// am I writing this wrong?
		log.Printf("key: %s  and value %s", pair.Key, pair.Value)
		parm.Add(pair.Key, pair.Value)
	}
	var lurl = fmt.Sprintf("https://api.lifx.com/v1/lights/id:%s/%s", bulbid, capability)
	req, err := http.NewRequest(http_method, lurl, nil)
	if err != nil {
		log.Println(err.Error())
	}
	req.URL.RawQuery = parm.Encode()
	log.Println(req.URL.RawQuery)
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

func (v *vislogger) priorityCheck(newPriority int) (string, string) {
	persist := "true"
	body := lnquery(v.bulbid, "", "GET", []parampair{})
	var data []interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	colour, err := dotpath.Get("color.hue", data[0], true)
	if v.colours[colour.(float64)].priority > newPriority {
		persist = "false"
	}
	return persist, fmt.Sprintf("%f", colour.(float64))
}

func (v *vislogger) Error() {
	persist, current_colour := v.priorityCheck(v.colours[0.0].priority)
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("hue:%f saturation:1.0", v.colours[0.0].hue)}, {Key: "period", Value: "2.0"}, {Key: "cycles", Value: "5.0"}, {Key: "from_color", Value: fmt.Sprintf("hue:%s", current_colour)}, {Key: "persist", Value: persist}}
	lnquery(v.bulbid, "effects/pulse", "POST", pairs)
}

func (v *vislogger) Warning() {
	persist, current_colour := v.priorityCheck(v.colours[36.0].priority)
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("hue:%f saturation:1.0", v.colours[36.0].hue)}, {Key: "period", Value: "2.0"}, {Key: "cycles", Value: "5.0"}, {Key: "from_color", Value: fmt.Sprintf("hue:%s", current_colour)}, {Key: "persist", Value: persist}}
	lnquery(v.bulbid, "effects/pulse", "POST", pairs)
}

func (v *vislogger) Normal() {
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("kelvin:%f  saturation:%f", v.base.kelvin, v.base.saturation)}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *vislogger) RGB(rgb string) {
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("rgb:%s saturation:1.0", rgb)}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *vislogger) HEX(hex string) {
	if !strings.HasPrefix(hex, "#") {
		hex = fmt.Sprintf("#%s", hex)
	}
	pairs := []parampair{{Key: "color", Value: fmt.Sprintf("%s saturation:1.0", hex)}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *vislogger) Party() {
	min := 0
	max := 255
	for {
		rgb := fmt.Sprintf("%d,%d,%d", rand.Intn(max-min)+min, rand.Intn(max-min)+min, rand.Intn(max-min)+min)
		pairs := []parampair{{Key: "color", Value: fmt.Sprintf("rgb:%s saturation:1.0", rgb)}}
		lnquery(v.bulbid, "state", "PUT", pairs)
		time.Sleep(500 * time.Millisecond)
	}
}

func (v *vislogger) Dim() {
	pairs := []parampair{{Key: "color", Value: "brightness:0.5"}}
	lnquery(v.bulbid, "state", "PUT", pairs)
}

func (v *vislogger) IsOn() bool {
	log.Println("Check if light is on")
	return true
}

func (v *vislogger) IsReady() bool {
	log.Println("check if light is ready")
	return true
}
