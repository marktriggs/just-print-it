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
	"strings"
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

	// Set some common paths to look for LibreOffice
	os.Setenv("PATH", os.Getenv("PATH")+":/Applications/LibreOffice.app/Contents/MacOS/")

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

func sanitizeFile(filename string) string {
	return filepath.Base(filename)
}

func convertToPDF(filename string) (string, error) {
	cmd := exec.Command("soffice", "--convert-to", "pdf",  "--outdir", filepath.Dir(filename), filename)
	ext := filepath.Ext(filename)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	if strings.Contains(string(output), "Error: source file could not be loaded") {
		return "", fmt.Errorf(string(output))
	}

	return filename[0:len(filename)-len(ext)] + ".pdf", nil
}

func convertIfNecessary(filename string) (string, error) {
	switch filepath.Ext(strings.ToLower(filename)) {
	case ".pdf", ".txt", ".ps", "":
		return filename, nil
	default:
		return convertToPDF(filename)
	}
}

func handleUpload(config *appConfig, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("upload")

	if err != nil {
		redirectWithStatus("upload_failed", w, r)
		return
	}

	outdir, err := ioutil.TempDir(os.TempDir(), "printme")

	if err != nil {
		redirectWithStatus("temp_dir_failed", w, r)
		return
	}

	defer os.RemoveAll(outdir)

	tempfilePath := path.Join(outdir, sanitizeFile(header.Filename))

	outfile, err := os.Create(tempfilePath)

	if _, err = io.Copy(outfile, file); err != nil {
		redirectWithStatus("temp_file_failed", w, r)
		return
	}

	outfile.Close()

	targetPath, err := convertIfNecessary(tempfilePath)

	if err != nil {
		redirectWithStatus("conversion_failed", w, r)
		fmt.Println(err)
		return
	}

	// Use this parameter to ensure our output doesn't get truncated...
	cmd := exec.Command("lpr", "-o", "fit-to-page", "-P", config.Printer, targetPath)
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
