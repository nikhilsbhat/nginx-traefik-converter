package middleware

import (
	"strings"

	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/configs"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	traefik "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/* ---------------- CONFIGURATION SNIPPET ---------------- */

// ConfigurationSnippet handles the below annotations.
// Annotations:
//   - "nginx.ingress.kubernetes.io/configuration-snippet"
func ConfigurationSnippet(ctx configs.Context) {
	ctx.Log.Debug("running converter ConfigurationSnippet")

	snippet, ok := ctx.Annotations["nginx.ingress.kubernetes.io/configuration-snippet"]
	if !ok || strings.TrimSpace(snippet) == "" {
		return
	}

	reqHeaders, respHeaders, warnings, unsupported := parseConfigurationSnippet(snippet)

	// Emit Warnings (gzip, cache, etc.)
	ctx.Result.Warnings = append(ctx.Result.Warnings, warnings...)

	// If there are unsupported directives (rewrite, lua, etc), do NOT convert
	if len(unsupported) > 0 {
		ctx.Result.Warnings = append(ctx.Result.Warnings,
			"configuration-snippet contains unsupported NGINX directives and was skipped",
		)

		return
	}

	// Nothing convertible
	if len(reqHeaders) == 0 && len(respHeaders) == 0 {
		return
	}

	// Create Headers middleware
	middleware := &traefik.Middleware{
		TypeMeta: metav1.TypeMeta{
			APIVersion: traefik.SchemeGroupVersion.String(),
			Kind:       "Middleware",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mwName(ctx, "snippet-headers"),
			Namespace: ctx.Namespace,
		},
		Spec: traefik.MiddlewareSpec{
			Headers: &dynamic.Headers{
				CustomRequestHeaders:  reqHeaders,
				CustomResponseHeaders: respHeaders,
			},
		},
	}

	ctx.Result.Middlewares = append(ctx.Result.Middlewares, middleware)
}

func parseConfigurationSnippet(snippet string) (map[string]string, map[string]string, []string, []string) {
	reqHeaders := map[string]string{}
	respHeaders := map[string]string{}
	warnings := make([]string, 0)
	unsupported := make([]string, 0)

	lines := strings.Split(snippet, "\n")

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		const proxySetHeaderCount = 3

		switch {
		case strings.HasPrefix(line, "more_set_headers"), // ───── Headers (convertible) ─────
			strings.HasPrefix(line, "add_header"):
			if h := extractHeader(line); h != nil {
				respHeaders[h[0]] = h[1]
			}
		case strings.HasPrefix(line, "proxy_set_header"): // proxy_set_header X-Foo bar;
			parts := strings.Fields(line)
			if len(parts) >= proxySetHeaderCount {
				reqHeaders[parts[1]] = strings.TrimSuffix(parts[2], ";")
			}
		case strings.HasPrefix(line, "gzip "): // ───── gzip (global-only in Traefik) ─────
			warnings = append(warnings,
				"gzip must be enabled globally in Traefik static configuration",
			)
		case strings.HasPrefix(line, "gzip_comp_level"):
			warnings = append(warnings,
				"gzip_comp_level is not configurable in Traefik and was ignored, compression level is fixed",
			)
		case strings.HasPrefix(line, "gzip_types"):
			warnings = append(warnings,
				"gzip_types is not configurable in Traefik and was ignored. Compresses a fixed, internal set of MIME types",
			)
		case strings.HasPrefix(line, "proxy_cache"): // ───── proxy_cache (not supported) ─────
			warnings = append(warnings,
				"proxy_cache is not supported in Traefik OSS and was ignored",
			)
		default: // ───── Everything else is unsafe ─────
			unsupported = append(unsupported, line)
		}
	}

	return reqHeaders, respHeaders, warnings, unsupported
}

func extractHeader(line string) []string {
	// expects: "X-Foo: bar"
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")

	if start == -1 || end <= start {
		return nil
	}

	const extractHeaderCount = 2

	keyValue := strings.SplitN(line[start+1:end], ":", extractHeaderCount)
	if len(keyValue) != extractHeaderCount {
		return nil
	}

	return []string{
		strings.TrimSpace(keyValue[0]),
		strings.TrimSpace(keyValue[1]),
	}
}
