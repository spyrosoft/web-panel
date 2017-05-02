package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type PanelConfig struct {
	ServerPort       string              `json:"server-port"`
	SiteDetails      []map[string]string `json:"site-detials"`
	NewSitePasswords []string            `json:"new-site-passwords"`
	Password         []byte              `json:"password"`
	PasswordSalt     []byte              `json:"password-salt"`
}

const (
	webRoot             = "awestruct/_site"
	panelConfigFilePath = "private/panel-config.json"
)

var (
	panelConfig = PanelConfig{}
	authToken   string
)

type apiResponse struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Fields   []string `json:"fields,omitempty"`
}

func main() {
	loadPanelConfig()
	router := httprouter.New()
	// Allows requests to pass through to NotFound if one method
	// is there and the other is not
	router.HandleMethodNotAllowed = false

	router.GET("/", authorize(serveStaticFilesOr404Handler))

	router.GET("/login/", redirectToHomeIfLoggedIn)
	router.GET("/logout/", logOut)

	router.POST("/api-noauth/", apiNoauth)
	router.POST("/api/", authorizeApi(api))

	router.NotFound = http.HandlerFunc(serveStaticFilesOr404)
	log.Fatal(http.ListenAndServe(":"+panelConfig.ServerPort, router))
}

func apiNoauth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	allowedFunctions := map[string]func(http.ResponseWriter, *http.Request) string{
		"login":  apiLogIn,
		"logout": apiLogOut,
	}
	apiGeneral(w, r, allowedFunctions)
}

func api(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	allowedFunctions := map[string]func(http.ResponseWriter, *http.Request) string{
		"update-site-file": apiUpdateSiteFile,
		"new-site":         apiNewSite,
	}
	apiGeneral(w, r, allowedFunctions)
}

func apiGeneral(w http.ResponseWriter, r *http.Request, allowedFunctions map[string]func(http.ResponseWriter, *http.Request) string) {
	if r.PostFormValue("action") == "" {
		serve404(w)
		return
	}
	function, ok := allowedFunctions[r.PostFormValue("action")]
	if !ok {
		json.NewEncoder(w).Encode(apiResponse{
			Errors: []string{"The requested action could not be found in our api: '" + r.PostFormValue("action") + "'"},
		})
		return
	}
	fmt.Fprint(w, function(w, r))
}

func serveStaticFilesOr404Handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serveStaticFilesOr404(w, r)
}

func (r apiResponse) String() string {
	resultsBytes, _ := json.Marshal(r)
	return string(resultsBytes)
}

func savePanelConfig() {
	panelConfigBytes, err := json.Marshal(panelConfig)
	panicOnErr(err)
	filePermissions := os.FileMode(0640)
	err = ioutil.WriteFile(panelConfigFilePath, panelConfigBytes, filePermissions)
	panicOnErr(err)
}

func debug(things ...interface{}) {
	fmt.Println("====================")
	for _, thing := range things {
		fmt.Printf("%+v\n", thing)
	}
	fmt.Println("====================")
}
