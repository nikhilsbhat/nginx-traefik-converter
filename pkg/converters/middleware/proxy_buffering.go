package middleware

import (
	"fmt"
	"strings"

	"github.com/nikhilsbhat/nginx-traefik-converter/pkg/configs"
	"github.com/nikhilsbhat/nginx-traefik-converter/pkg/converters/models"
)

func ProxyBuffering(ctx configs.Context) {
	ctx.Log.Debug("running converter ProxyBuffering")

	ann := string(models.ProxyBuffering)

	val, ok := ctx.Annotations[ann]
	if !ok {
		return
	}

	v := strings.ToLower(strings.TrimSpace(val))

	switch v {
	case "on":
		warningMessage := "nginx.ingress.kubernetes.io/proxy-buffering is not supported in Traefik and was ignored"

		ctx.Result.Warnings = append(
			ctx.Result.Warnings,
			"nginx.ingress.kubernetes.io/proxy-buffering is not supported in Traefik and was ignored",
		)

		ctx.ReportIgnored(ann, warningMessage)
	case "off":
		warningMessage := "proxy-buffering=off is default behavior in Traefik"

		ctx.Result.Warnings = append(ctx.Result.Warnings, warningMessage)

		ctx.ReportIgnored(ann, warningMessage)
	default:
		warningMessage := fmt.Sprintf(
			"nginx.ingress.kubernetes.io/proxy-buffering has unknown value %q and was ignored", val)

		ctx.Result.Warnings = append(ctx.Result.Warnings, warningMessage)

		ctx.ReportIgnored(ann, warningMessage)
	}
}
