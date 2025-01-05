package handlers

import (
	"context"
	"net/http"

	whweb "github.com/bohdanch-w/wheel/web"
	whapi "github.com/bohdanch-w/wheel/web/api"
)

func Version(version, build, buildTime string) whapi.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		resp := versionResponse{
			Version:   version,
			Build:     build,
			BuildTime: buildTime,
		}

		return whweb.Respond(w, http.StatusOK, resp)
	}
}

type versionResponse struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	BuildTime string `json:"build_time"`
}
