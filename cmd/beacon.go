package main

import (
	"log"
	"time"

	"github.com/wissam/beacon/pkg/vislog"
)

//fuck the tui, useless for now.. make a webhook experiment?

func main() {
	blb := vislog.NewBulb("d073d567639b")
	//blb.ShowAll()
	//blb.Normal()
	//blb.ShowAll()
	//j	blb.Normal()
	//	time.Sleep(15 * time.Second)
	blb.Warning() // should stay orange after this...
	//time.Sleep(5 * time.Second)
	time.Sleep(15 * time.Second)
	//blb.Warning()
	blb.Error() //should stay red
	//blb.Normal()
	time.Sleep(15 * time.Second)
	//logic is still wrong... should be false not true ...! brb wc
	blb.Warning() //should stay red
	//snssend.Send()
	//emailsend.Send()
	log.Println(blb.Status())
	//something wrong... my logic is flawed...
}
