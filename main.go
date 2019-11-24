package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	muxer := http.NewServeMux()
	muxer.HandleFunc("/upload", handleUpload)
	http.ListenAndServe(":8080", muxer)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	uploadedFile, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}
	defer uploadedFile.Close()

	arch := r.Form.Get("arch")
	repo := r.Form.Get("repo")
	token := r.Form.Get("token")
	if token != os.Getenv("PACKAGER_TOKEN") {
		w.WriteHeader(401)
		w.Write([]byte("Invalid authentification token"))
		return
	}

	fmt.Printf("Receiving file '%s'\n", header.Filename)

	if err = addPackage(arch, repo, header.Filename, uploadedFile); err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}
}

func addPackage(arch, repo, pkg string, data io.Reader) error {
	directoryPath := arch + "/" + repo
	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return err
	}

	packagePath := directoryPath + "/" + pkg
	packageFile, err := os.OpenFile(packagePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	written, err := io.Copy(packageFile, data)
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes into '%s'\n", written, packagePath)

	databasePath := directoryPath + "/" + repo + ".db.tar.gz"
	if err = updatePackageDatabase(databasePath, packagePath); err != nil {
		return err
	}

	return nil
}

func updatePackageDatabase(db, pkg string) error {
	fmt.Printf("Executing `repo-add %s %s`\n", db, pkg)
	cmd := exec.Command("repo-add", db, pkg)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s", output)

	return nil
}
