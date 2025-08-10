package main

import (
	"encoding/xml"
	"net/http"
)

type Dialplan struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Section struct {
		XMLName xml.Name `xml:"section"`
		Name    string   `xml:"name,attr"`
		Result  struct {
			XMLName xml.Name `xml:"result"`
			Status  string   `xml:"status,attr"`
			Node    struct {
				XMLName xml.Name `xml:"dialplan"`
				Context Context  `xml:"context"`
			} `xml:"data"`
		} `xml:"result"`
	} `xml:"section"`
}

type Context struct {
	Name  string  `xml:"name,attr"`
	Exten []Extension `xml:"extension"`
}

type Extension struct {
	Name    string    `xml:"name,attr"`
	Cond    []Condition `xml:"condition"`
}

type Condition struct {
	Field string `xml:"field,attr"`
	Expr  string `xml:"expression,attr"`
	Acts  []Action `xml:"action"`
}

type Action struct { Application string `xml:"application,attr"`; Data string `xml:"data,attr"` }

func handleDialplan(w http.ResponseWriter, r *http.Request) {
	ctx := Context{Name: "public"}
	// simple public: answer progress media then bridge to internal or hangup
	ext := Extension{Name: "outbound"}
	cond := Condition{Field: "destination_number", Expr: ".*"}
	cond.Acts = []Action{
		{Application: "export", Data: "ignore_early_media=false"},
		{Application: "set", Data: "hangup_after_bridge=true"},
		{Application: "set", Data: "record_session=/recordings/${uuid}.wav"},
	}
	ext.Cond = []Condition{cond}
	ctx.Exten = []Extension{ext}

	doc := Dialplan{Type: "freeswitch/xml"}
	doc.Section.Name = "dialplan"
	doc.Section.Result.Status = "success"
	doc.Section.Result.Node.Context = ctx

	w.Header().Set("Content-Type", "text/xml")
	xml.NewEncoder(w).Encode(doc)
}