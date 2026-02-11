package middleware

import (
	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/configs"
	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/converters/models"
)

/* ---------------- BODY SIZE ---------------- */

// EnableUnderscoresInHeaders handles the below annotations.
// Annotations:
//   - "nginx.ingress.kubernetes.io/enable-underscores-in-headers"
func EnableUnderscoresInHeaders(ctx configs.Context) {
	ann := string(models.UnderscoresInHeaders)

	if _, ok := ctx.Annotations[ann]; !ok {
		return
	}

	ctx.ReportIgnored(
		ann,
		"Traefik accepts headers with underscores by default (it uses Go HTTP parser); no equivalent configuration is required",
	)
}
