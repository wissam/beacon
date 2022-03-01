package hooks

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wissam/beacon/pkg/vislog"
)

type ghhook struct {
	Action     string `json:"action"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Workflow_Run struct {
		Id         int    `json:"id"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	} `json:"workflow_run"`
}

//now to make it proper...
func handleGHhook(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	ghh := ghhook{}
	json.Unmarshal(b, &ghh)
	log.Printf("Repo name %s\n", ghh.Repository.Name)
	log.Printf("Workflow Run id %d\n", ghh.Workflow_Run.Id)
	log.Printf("Status: %s\n", ghh.Workflow_Run.Status)
	log.Printf("Conclusion: %s\n", ghh.Workflow_Run.Conclusion)
	// I didn't even run the function here lol...
	if ghh.Workflow_Run.Conclusion != "" {
		lightMeBulb("d073d567639b", ghh.Workflow_Run.Conclusion)
	}
}

func lightMeBulb(bulbid string, conclusion string) {
	blb := vislog.NewBulb(bulbid)
	m := map[string]func(){
		"failure":   blb.Error,
		"cancelled": blb.Warning,
		"success":   blb.Success,
	}
	val := m[conclusion]
	val()
}

//this keeps happening~~~
func Run() {
	log.Println("Hooks Server Started!")
	http.HandleFunc("/ghhook", handleGHhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
