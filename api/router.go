package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

const host = "newapi.watermelonpig.eu.org:17728"

func init() {
	router = gin.Default()
	router.Any("/*path", func(ctx *gin.Context) {
		ctx.Request.URL.Host = host
		ctx.Request.URL.Scheme = "https"
		r, err := http.NewRequest(ctx.Request.Method, ctx.Request.URL.String(), ctx.Request.Body)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		r.Header = ctx.Request.Header
		r.Header.Set("Host", host)
		r.Header.Set("X-Forwarded-Proto", ctx.Request.URL.Scheme)
		r.Header.Set("X-Forwarded-For", strings.Join(ctx.Request.Header["X-Forwarded-For"], ",")+","+ctx.Request.RemoteAddr)
		r.Header.Set("X-Real-IP", ctx.Request.RemoteAddr)
		r.Header.Del("Accept-Encoding")
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer resp.Body.Close()
		ctx.Stream(func(w io.Writer) bool {
			i, err := io.Copy(w, resp.Body)
			return i == 0 || err != nil
		})
	})
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
