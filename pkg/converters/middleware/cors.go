package middleware

import (
	"strconv"
	"strings"

	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/configs"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	traefik "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/* ---------------- CORS ---------------- */

// CORS handles the below annotations.
// Annotations:
//   - "nginx.ingress.kubernetes.io/enable-cors"
//   - "nginx.ingress.kubernetes.io/cors-allow-origin"
//   - "nginx.ingress.kubernetes.io/cors-allow-methods"
//   - "nginx.ingress.kubernetes.io/cors-allow-headers"
//   - "nginx.ingress.kubernetes.io/cors-allow-credentials"
//   - "nginx.ingress.kubernetes.io/cors-max-age"
//   - "nginx.ingress.kubernetes.io/cors-expose-headers"
//   - "nginx.ingress.kubernetes.io/cors-expose-headers"
func CORS(ctx configs.Context) error {
	ctx.Log.Debug("running converter CORS")

	if ctx.Annotations["nginx.ingress.kubernetes.io/enable-cors"] != "true" {
		return nil
	}

	headers := &dynamic.Headers{}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-allow-origin"]; v != "" {
		headers.AccessControlAllowOriginList = headersNeat(v)
	}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-allow-methods"]; v != "" {
		headers.AccessControlAllowMethods = headersNeat(v)
	}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-allow-headers"]; v != "" {
		headers.AccessControlAllowHeaders = headersNeat(v)
	}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-allow-credentials"]; v == "true" {
		headers.AccessControlAllowCredentials = true
	}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-max-age"]; v != "" {
		secs, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}

		headers.AccessControlMaxAge = secs
	}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-expose-headers"]; v != "" {
		headers.AccessControlExposeHeaders = headersNeat(v)
	}

	if v := ctx.Annotations["nginx.ingress.kubernetes.io/cors-expose-headers"]; v != "" {
		headers.AccessControlExposeHeaders = headersNeat(v)
	}

	ctx.Result.Middlewares = append(ctx.Result.Middlewares, &traefik.Middleware{
		TypeMeta: metav1.TypeMeta{
			APIVersion: traefik.SchemeGroupVersion.String(),
			Kind:       "Middleware",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mwName(ctx, "cors"),
			Namespace: ctx.Namespace,
		},
		Spec: traefik.MiddlewareSpec{Headers: headers},
	})

	return nil
}

func headersNeat(value string) []string {
	headers := strings.Split(value, ",")

	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	return headers
}
