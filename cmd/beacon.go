package main

import (
	"github.com/wissam/beacon/internal/vislog"
)

func main() {
	blb := vislog.NewBulb("d073d567639b")
	//log.Println("Pulse from white to red 5 times, persists as red")
	//blb.Error()
	//time.Sleep(12 * time.Second)
	//log.Println("Pulse from red to orange 5 times, reverts back to red")
	//blb.Warning()
	//time.Sleep(12 * time.Second)
	//log.Println("Goes back to white")
	//blb.Normal()
	//SUCCESS~
	// now what? :P
	// I don't like how the colours are organised still, it is the wrong
	// datastructure , so 1 I either go search for the right one and find it
	// somehow? or go thru the webserver work to build webhooks and test it
	// with github
	//blb.RGB("255,0,255")
	//yay success...
	//blb.Normal()
	//	blb.Dim()
	blb.Party()
}
