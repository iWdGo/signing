package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// sentence if fixed as root page is tested
var want = "Hi there, I love Go!"

// TODO add login/logout tests
func TestMyHandlerOffline(t *testing.T) {
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal("New request failed with ", err)
	}
	defer r.Body.Close()

	w := httptest.NewRecorder()
	homeHandler(w, r)

	if w.Code != 200 {
		t.Fatalf("wrong code returned: %d", w.Code)
	}

	got := w.Body.String()
	if strings.Index(string(got), want) != 0 { // ignoring footer
		t.Fatalf("wrong body returned: %s", got)
	}
}

/* Panic-ing when closing Body and no response otherwise */
func TestHandlerOnLine(t *testing.T) {
	client := &http.Client{}

	r, err := client.Get(url) // "/" is the same page
	if err != nil {
		t.Fatal(url, "is unavailable: ", err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		t.Fatal("client get is not OK: ", r.StatusCode)
	}

	got := make([]byte, r.ContentLength)
	b, err := r.Body.Read(got)
	if int64(b) != r.ContentLength {
		t.Fatal("data lost")
	}
	if err != io.EOF {
		t.Fatal("error reading body: ", err, "and read", b)
	} else if strings.Index(string(got), want) != 0 { // ignoring footer
		t.Fatalf("wrong body returned: %s", got)
	}
}
