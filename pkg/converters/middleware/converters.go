package middleware

import (
	"strings"

	"github.com/nikhilsbhat/nginx-traefik-converter/pkg/configs"
)

/* ---------------- WARNINGS ---------------- */

// Warnings adds warnings to the parsed annotations if any.
func Warnings(ctx configs.Context) {
	for annotation := range ctx.Annotations {
		if strings.Contains(annotation, "auth-tls") ||
			strings.Contains(annotation, "snippet") ||
			strings.Contains(annotation, "proxy-read") ||
			strings.Contains(annotation, "proxy-send") {
			ctx.Result.Warnings = append(ctx.Result.Warnings, annotation+" is not safely convertible")
		}
	}
}

func mwName(ctx configs.Context, suffix string) string {
	return ctx.IngressName + "-" + suffix
}
