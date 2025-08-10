package main

import (
	"encoding/xml"
	"net/http"
	"os"
)

type Directory struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Section struct {
		XMLName xml.Name `xml:"section"`
		Name    string   `xml:"name,attr"`
		Result  struct {
			XMLName xml.Name `xml:"result"`
			Status  string   `xml:"status,attr"`
			Node    struct {
				XMLName xml.Name `xml:"directory"`
				Domain  struct {
					XMLName xml.Name `xml:"domain"`
					Name    string   `xml:"name,attr"`
					Params  []Param  `xml:"params>param"`
					Users   []User   `xml:"groups>group>users>user"`
				} `xml:"domain"`
			} `xml:"data"`
		} `xml:"result"`
	} `xml:"section"`
}

type Param struct{ Name string `xml:"name,attr"`; Value string `xml:"value,attr"` }

type User struct{ ID string `xml:"id,attr"`; Params []Param `xml:"params>param"` }

func handleDirectory(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" { domain = "default" }
	doc := Directory{Type: "freeswitch/xml"}
	doc.Section.Name = "directory"
	doc.Section.Result.Status = "success"
	doc.Section.Result.Node.Domain.Name = domain
	doc.Section.Result.Node.Domain.Params = []Param{{Name: "dial-string", Value: "{sip_invite_params=alert_info=ring;info}sofia/internal/$${destination_number}"}}
	w.Header().Set("Content-Type", "text/xml")
	xml.NewEncoder(w).Encode(doc)
}

func main() {
	addr := ":9000"
	if v := os.Getenv("CFG_HTTP_ADDR"); v != "" { addr = v }
	http.HandleFunc("/directory", handleDirectory)
	http.ListenAndServe(addr, nil)
}