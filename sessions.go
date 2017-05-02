package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func isLoggedIn(w http.ResponseWriter, r *http.Request) (isLoggedIn bool) {
	cookieToken, ok := cookieStringValue("auth-token", r)
	if !ok {
		return
	}
	if authToken == "" {
		return
	}
	if authToken != cookieToken {
		logOutUser(w, r)
		return
	}
	return true
}

func cookieIntValue(cookieName string, r *http.Request) (value int, ok bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	value, err = strconv.Atoi(cookie.Value)
	if err != nil {
		return
	}
	ok = true
	return
}

func cookieStringValue(cookieName string, r *http.Request) (value string, ok bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	value = cookie.Value
	ok = true
	return
}

func apiLogIn(w http.ResponseWriter, r *http.Request) string {
	if isLoggedIn(w, r) {
		return apiResponse{Success: true}.String()
	}
	err := logInUser(r.PostFormValue("password"), w)
	if err != nil {
		return apiResponse{
			Messages: []string{err.Error()},
			Fields:   []string{"password"},
		}.String()
	}
	return apiResponse{Success: true}.String()
}

func logInUser(password string, w http.ResponseWriter) (err error) {
	ok, err := checkPassword(password)
	if err != nil {
		return
	}
	if !ok {
		return errors.New("Incorrect password.")
	}
	authToken, err = generateStringToken(30)
	if err != nil {
		return
	}
	newCookie := http.Cookie{Name: "auth-token", Value: authToken, Path: "/"}
	http.SetCookie(w, &newCookie)
	return
}

func checkPassword(password string) (ok bool, err error) {
	if len(panelConfig.Password) == 0 && len(panelConfig.PasswordSalt) == 0 {
		return signUp(password)
	}
	hash, err := scryptHash(password, panelConfig.PasswordSalt)
	if err != nil {
		return
	}
	if string(hash) == string(panelConfig.Password) {
		ok = true
	}
	return
}

func signUp(password string) (ok bool, err error) {
	ok, err = newPanelConfigPasswords(password)
	if !ok {
		return
	}
	savePanelConfig()

	ok = true
	return
}

func newPanelConfigPasswords(password string) (ok bool, err error) {
	var hash, salt []byte
	hash, salt, err = scryptHashAndSalt(password)
	if err != nil {
		return
	}
	panelConfig.Password = hash
	panelConfig.PasswordSalt = salt

	for i := 0; i < 15; i++ {
		var newSitePassword string
		newSitePassword, err = generateStringToken(50)
		if err != nil {
			return
		}
		panelConfig.NewSitePasswords = append(
			panelConfig.NewSitePasswords,
			newSitePassword,
		)
	}

	err = errors.New("Since this is the first time you are setting up the site, the one-time passwords to create new sites have been generated for the first time. Please store these somewhere safe: " + fmt.Sprintf("%+v\n", panelConfig.NewSitePasswords))

	return
}

func redirectToHomeIfLoggedIn(w http.ResponseWriter, r *http.Request, requestParameters httprouter.Params) {
	if isLoggedIn(w, r) {
		http.Redirect(w, r, "/", 302)
		return
	}
	serveStaticFilesOr404(w, r)
}

func logOut(w http.ResponseWriter, r *http.Request, requestParameters httprouter.Params) {
	logOutUser(w, r)
	http.Redirect(w, r, "/login/", 302)
}

func apiLogOut(w http.ResponseWriter, r *http.Request) (results string) {
	logOutUser(w, r)
	results = apiResponse{Success: true}.String()
	return
}

func logOutUser(w http.ResponseWriter, r *http.Request) {
	authToken = ""
	deleteCookie := http.Cookie{Name: "user-id", MaxAge: -1}
	http.SetCookie(w, &deleteCookie)
	deleteCookie = http.Cookie{Name: "auth-token", MaxAge: -1}
	http.SetCookie(w, &deleteCookie)
}

func serveLoginPage(w http.ResponseWriter) {
	template, err := ioutil.ReadFile(webRoot + "/login/index.html")
	if err != nil {
		serve500(w)
	}
	fmt.Fprint(w, string(template))
}
