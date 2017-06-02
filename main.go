package main

/*
 * ./restad2pve-groups --dest-group=CSR --pve-relm=byu --src-group=physics-csrs --restad-server=avari:1234> users.cfg
 *  N.B., this depends on RESTAD (see http://github.com/kalebo/RESTAD ) being running
 */

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"text/template"
	"time"
)

type GroupUsers struct {
	GroupName string
	Relm      string
	Users     []RestADUser
}

type RestADUser struct {
	Name  string
	NetId string
}

type RestADMembers struct {
	Name  string
	Users []RestADUser
}

const tpml_str string = `
user:root@pam:1:0:::physics-csr@byu.edu:::
{{range .Users}}user:{{.NetId}}@{{$.Relm}}:1:0:{{.Name}}:::::
{{else}}{{end}}

group:{{.GroupName}}:{{range $index, $element := .Users}}{{if $index}},{{end}}{{.NetId}}@{{$.Relm}}{{else}}{{end}}::


acl:1:/:@{{.GroupName}}:Administrator:
`

func main() {
	destPtr := flag.String("dest-group", "", "The PVE group to have administrator role and users added")
	srcPtr := flag.String("src-group", "", "The AD group to pull the users from")
	relmPtr := flag.String("pve-relm", "", "The realm in pve that corresponds to your AD")
	restadPtr := flag.String("restad-server", "", "The hostname and port of your RESTAD server connected to AD")
	flag.Parse()

	if *destPtr == "" || *srcPtr == "" || *relmPtr == "" || *restadPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	http_client := &http.Client{Timeout: 5 * time.Second} // make sure we don't wait forever

	resp, err := http_client.Get("http://" + *restadPtr + "/api/group/" + *srcPtr + "/members")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var members RestADMembers

	json.NewDecoder(resp.Body).Decode(&members)

	tmpl, _ := template.New("user.cfg").Parse(tpml_str)

	params := GroupUsers{*destPtr, *relmPtr, members.Users}

	tmpl.Execute(os.Stdout, params)
}
