package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type appConfig struct {
	BaseDir string
	Printer string
}

type appHandler struct {
	config  *appConfig
	handler func(*appConfig, http.ResponseWriter, *http.Request)
}

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler(h.config, w, r)
}

func main() {
	basedir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port> <printer name>\n", os.Args[0])
		return
	}

	port := ":" + os.Args[1]
	printer := os.Args[2]

	fmt.Printf("Listening on port %s.  Will print to printer '%s'\n", port, printer)

	config := &appConfig{basedir, printer}

	http.Handle("/", appHandler{config, showIndex})
	http.Handle("/upload", appHandler{config, handleUpload})

	http.ListenAndServe(port, nil)
}

func showIndex(config *appConfig, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	template, err := template.ParseFiles(path.Join(config.BaseDir, "templates", "index.html"))

	if err != nil {
		log.Fatal(err)
	}

	template.Execute(w, struct {
		Status string
	}{
		r.FormValue("status"),
	})
}

func redirectWithStatus(status string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?status="+status, 302)
}

func handleUpload(config *appConfig, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("upload")

	if err != nil {
		redirectWithStatus("upload_failed", w, r)
		return
	}

	outfile, err := ioutil.TempFile(os.TempDir(), "printme")

	if err != nil {
		redirectWithStatus("temp_file_failed", w, r)
		return
	}

	defer os.Remove(outfile.Name())

	if _, err = io.Copy(outfile, file); err != nil {
		redirectWithStatus("temp_file_failed", w, r)
		return
	}

	cmd := exec.Command("lpr", "-P", config.Printer, outfile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		fmt.Println(string(output))
		redirectWithStatus("print_failed", w, r)
		return
	}

	redirectWithStatus("ok", w, r)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
