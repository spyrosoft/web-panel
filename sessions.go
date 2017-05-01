package main

import (
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
	authToken, err := generateStringToken(30)
	if err != nil {
		return
	}
	newCookie := http.Cookie{Name: "auth-token", Value: authToken, Path: "/"}
	http.SetCookie(w, &newCookie)
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
