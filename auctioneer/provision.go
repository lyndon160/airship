package auctioneer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type Provision struct {
	Nodes        []string       `json:"nodes"`
	ImageName    string         `json:"image_name"`
	Memory       int            `json:"ram"`
	Hours        int            `json:"hours"`
	PortBindings map[string]int `json:"port_bindings"`
	UserID       string         `json:"user_id,omitempty"`
	AuctionID    string         `json:"auction_id,omitempty"`
}

//[List nodes] [String image_name] [Int ram] [Int hours]
//{"image_name": "hypriot/rpi-busybox-httpd", "nodes": ["192.168.2.15","192.168.2.18"],"port_bindings": {"internal": 80, "external": 64444}, "ram":"200", "hours": 24}

func checkValidity(col *mgo.Collection, auctionID string) bool {
	// query := bson.M{"id": auctionID}

	// auction := Auction{}
	// if e := col.Find(query).One(&auction); e != nil {
	// 	log.Panic(e)
	// }

	//check if auction is not live any more

	return true
}

func resolveNodes() []string {
	var nodes []string
	return nodes
}

func provision(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var provision Provision

	body := readRequestBody(r)

	if e := json.Unmarshal(body, &provision); e != nil {
		log.Panic(e)
	}

	col := a.session.DB(configuration.DatabaseName).C(collectionNames["auction"])

	if checkValidity(col, provision.AuctionID) {
		// provision.Nodes = resolveNodes()
		provision.UserID = ""
		provision.AuctionID = ""

		post, err := json.Marshal(provision)
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(post)

		path := configuration.ProvisionerPath + "/nodes/provision_dockers"
		log.Println(path)
		resp, err := http.Post(path, "application/json; charset=UTF-8", reader) //should the end point not be '/provision_docker_containers'?
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		log.Println(string(body))
	}

	return 200, nil
}
