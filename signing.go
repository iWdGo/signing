package main

// https://cloud.google.com/appengine/docs/standard/go/users/#Go_Google_accounts_and_the_development_server
import (
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"log"
	"net/http"
)

const (
	address = "localhost:8080"
	url     = "http://" + address
)

var (
	luser = ""
)

/* Navigation */
func footer(w http.ResponseWriter) {
	fmt.Fprintf(w, `<a href="%s">Home</a>  `, url)
	if luser == "" {
		fmt.Fprintf(w, `<a href="%s">Log in</a><br>`, url+"/login")
	} else {
		fmt.Fprintf(w, `<a href="%s">Log out</a><br>`, url+"/logout")
	}
}

/* Login handler using standard user */
func userLogin(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	u := user.Current(c) // displays the login form
	if u == nil {
		loginUrl, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// You can start standard login by setting a 302 Status code on the header
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusFound)
		// Or you request a confirmation
		// fmt.Fprintf(w, `<a href="%s">Log in</a>`, loginUrl)
	} else {
		if luser == "" { // context is lost after restart. So always greeted once.
			luser = u.Email
			fmt.Fprintf(w, `Hello, %v!<br>`, luser)
		} else {
			fmt.Fprintf(w, `Already logged in, %v!<br>`, luser)
		}
		footer(w)
	}
}

/* Logout */
func userLogout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	u := user.Current(c)
	if u == nil {
		// You may not write on the request if header is set
		if l := w.Header().Get("Location"); l != "" {
			log.Println("Location =", l)
		} else {
			fmt.Fprintf(w, "Nobody logged in<br>")
		}
	} else {
		logoutURL, err := user.LogoutURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", logoutURL) // to execute the logout
		// u is nil now
		fmt.Fprintf(w, "Good bye, %s!<br>", u.Email)
		// w.WriteHeader(http.StatusFound) // Logout is automatic, i.e. the following will not print
		fmt.Fprintf(w, `Please confirm by clicking <a href="%s">Log out</a><br>`, logoutURL)
		// logout is actually a 302 request
		// GET /_ah/logout?continue=http%3A//localhost%3A8080/logout HTTP/1.1" 302
		// if w.WriteHeader(http.StatusFound), you can't say goodbye
		luser = "" // User is logged out
	}
	footer(w) // Print footer before erasing logged out user.
	// Log in will be the cancel as login check context and not the variable
}

/* Home page */
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// when starting using dev_appserver.py, a request on / is sent during start up
	w.Header().Set("Content-type", "text/html; charset=utf-8") // Fprint is html
	puser := luser
	s := "s" // 3rd
	if puser == "" {
		puser = "I"
		s = "" // 1st
	}
	fmt.Fprintf(w, "Hi there, %s love%s Go!<br>", puser, s)
	footer(w)
}

/* main() means go run or flex deployment */
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", userLogin)
	http.HandleFunc("/logout", userLogout)
	if appengine.IsDevAppServer() {
		appengine.Main()
	} else {
		// deployed using go run
		http.ListenAndServe(address, nil)
	}
}
