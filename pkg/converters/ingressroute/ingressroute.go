package ingressroute

import (
	"fmt"

	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/configs"
	"github.com/nikhilsbhat/ingress-traefik-converter/pkg/converters/tls"
	traefik "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// BuildIngressRoute handles the below annotations.
// Annotations:
//   - "nginx.ingress.kubernetes.io/backend-protocol"
//   - "nginx.ingress.kubernetes.io/grpc-backend"
func BuildIngressRoute(ctx configs.Context) error {
	ing := ctx.Ingress

	// 1️⃣ Resolve backend protocol ONCE (Ingress-wide)
	scheme, err := resolveScheme(ctx.Annotations)
	if err != nil {
		return err
	}

	routes := make([]traefik.Route, 0)
	seen := make(map[string]struct{}) // dedup key set

	for _, rule := range ing.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}

		hostMatch := buildHostMatch(rule.Host)

		for _, path := range rule.HTTP.Paths {
			svc := path.Backend.Service
			if svc == nil {
				continue
			}

			pathMatch := buildPathMatch(path)
			match := combineMatch(hostMatch, pathMatch)

			// Build a stable dedup key
			key := fmt.Sprintf(
				"host=%s|path=%s|svc=%s|port=%d|scheme=%s",
				rule.Host,
				path.Path,
				svc.Name,
				svc.Port.Number,
				scheme,
			)

			if _, exists := seen[key]; exists {
				continue // skip duplicate route
			}

			seen[key] = struct{}{}

			route := traefik.Route{
				Kind:  "Rule",
				Match: match,
				Services: []traefik.Service{
					{
						LoadBalancerSpec: traefik.LoadBalancerSpec{
							Name: svc.Name,
							Port: intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: svc.Port.Number,
							},
							Scheme: scheme,
						},
					},
				},
				Middlewares: middlewareRefs(ctx),
			}

			routes = append(routes, route)
		}
	}

	if len(routes) == 0 {
		return nil
	}

	ingressRoute := &traefik.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: traefik.SchemeGroupVersion.String(),
			Kind:       "IngressRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ing.Name,
			Namespace: ing.Namespace,
		},
		Spec: traefik.IngressRouteSpec{
			EntryPoints: entryPointsForScheme(scheme),
			Routes:      routes,
		},
	}

	// Apply TLS only if scheme requires it (as discussed earlier)
	tls.ApplyTLSOption(ingressRoute, ctx, scheme)

	ctx.Result.IngressRoutes = append(ctx.Result.IngressRoutes, ingressRoute)

	return nil
}

func middlewareRefs(ctx configs.Context) []traefik.MiddlewareRef {
	refs := make([]traefik.MiddlewareRef, 0)

	for _, mw := range ctx.Result.Middlewares {
		refs = append(refs, traefik.MiddlewareRef{
			Name: mw.GetName(),
		})
	}

	return refs
}

func buildHostMatch(host string) string {
	if host == "" {
		return ""
	}

	return fmt.Sprintf("Host(`%s`)", host)
}

func buildPathMatch(path netv1.HTTPIngressPath) string {
	pth := path.Path
	if pth == "" {
		pth = "/"
	}

	switch *path.PathType {
	case netv1.PathTypeExact:
		return fmt.Sprintf("Path(`%s`)", pth)
	case netv1.PathTypePrefix:
		return fmt.Sprintf("PathPrefix(`%s`)", pth)
	case netv1.PathTypeImplementationSpecific:
		// Best-effort: treat as Prefix (same as most ingress controllers do)
		return fmt.Sprintf("PathPrefix(`%s`)", pth)
	default:
		return fmt.Sprintf("PathPrefix(`%s`)", pth)
	}
}

func combineMatch(hostMatch, pathMatch string) string {
	switch {
	case hostMatch != "" && pathMatch != "":
		return hostMatch + " && " + pathMatch
	case hostMatch != "":
		return hostMatch
	case pathMatch != "":
		return pathMatch
	default:
		return "PathPrefix(`/`)"
	}
}
