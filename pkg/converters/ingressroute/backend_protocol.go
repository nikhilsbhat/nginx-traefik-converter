package ingressroute

import (
	"strings"

	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/converters/models"
	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/errors"
)

func resolveScheme(annotations map[string]string) (string, error) {
	backendProto := strings.ToUpper(annotations[string(models.BackendProtocol)])
	grpcBackend := annotations[string(models.GrpcBackend)] == "true"

	if grpcBackend && backendProto == "HTTP" {
		return "h2c", nil // but you could also emit a warning via ctx
	}

	// If backend-protocol is explicitly set, it should take precedence
	switch backendProto {
	case "", "HTTP":
		if grpcBackend {
			// gRPC without TLS
			return "h2c", nil
		}

		return "http", nil

	case "HTTPS":
		return "https", nil

	case "GRPC":
		// gRPC without TLS
		return "h2c", nil

	case "GRPCS":
		// gRPC over TLS
		return "https", nil

	default:
		return "", &errors.ConverterError{Message: "unsupported backend-protocol"}
	}
}

func entryPointsForScheme(scheme string) []string {
	switch scheme {
	case "https":
		return []string{"websecure"}
	case "h2c":
		return []string{"web"}
	default:
		return []string{"web"}
	}
}

// NeedsIngressRoute makes the decision on requirement of ingress routes.
func NeedsIngressRoute(ann map[string]string) bool {
	if ann[string(models.GrpcBackend)] == "true" {
		return true
	}

	if _, ok := ann[string(models.BackendProtocol)]; ok {
		return true
	}

	return false
}
