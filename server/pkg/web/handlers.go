package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"istio-redirector/domain"
	"istio-redirector/pkg/github"
	"istio-redirector/pkg/redirections"
	"net/http"
	"strconv"

	"github.com/n0rad/go-erlog/logs"
)

func UploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("csv_file")
	if err != nil {
		logs.WithE(err).Error("can't get csv file from form")
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return
	}

	payload, err := redirections.Generate(
		domain.InputData{
			File:            buf.Bytes(),
			RedirectionName: r.FormValue("redirection_name"),
			RedirectionType: r.FormValue("redirection_type"),
		},
	)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	pushGithub, err := strconv.ParseBool(r.FormValue("pushGithub"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	if pushGithub {
		prURL := github.Create(payload.Bytes(), r.FormValue("redirection_name"), r.FormValue("redirection_env"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"PR": prURL})
		return
	}

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=virtual-service.yaml")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", payload.Len()))

	//stream the body to the client without fully loading it into memory
	io.Copy(w, &payload)
}
