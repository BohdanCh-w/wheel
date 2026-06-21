package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/exp/constraints"

	"github.com/bohdanch-w/wheel/web"
	"github.com/bohdanch-w/wheel/web/api"
)

type ResourceMid[T constraints.Integer | string | uuid.UUID] struct {
	Name string
}

func (mid *ResourceMid[T]) Wrap(h api.Handler) api.Handler {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var (
			etalon T
			err    error
			id     any
		)

		switch any(etalon).(type) {
		case string:
			id = mux.Vars(r)[mid.Name]
		case uint:
			id, err = parseInteger[uint](mux.Vars(r)[mid.Name])
		case int:
			id, err = parseInteger[int](mux.Vars(r)[mid.Name])
		case uuid.UUID:
			id, err = uuid.Parse(mux.Vars(r)[mid.Name])
		}

		if err != nil {
			return web.NewError(http.StatusNotFound, fmt.Errorf("parse %q: %w", mid.Name, err))
		}

		ctx = context.WithValue(ctx, pathParamKey(mid.Name), id)

		return h(ctx, w, r)
	}

	return f
}

func parseInteger[T constraints.Integer](s string) (T, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return T(0), err // nolint: wrapcheck
	}

	return T(value), nil
}

type pathParamKey string

func GetPathParam[T any](ctx context.Context, name string) T {
	value, _ := ctx.Value(pathParamKey(name)).(T)

	return value
}
