package configs

import (
	"log/slog"

	netv1 "k8s.io/api/networking/v1"
)

// Context holds the necessary info required to run the converters.
type Context struct {
	Ingress     *netv1.Ingress    `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	IngressName string            `yaml:"ingress_name,omitempty" json:"ingress_name,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Result      *Result           `yaml:"result,omitempty" json:"result,omitempty"`
	Log         *slog.Logger
}

// New returns a new instance of Context when invoked.
func New(ingress *netv1.Ingress, result *Result, logger *slog.Logger) *Context {
	return &Context{
		Ingress:     ingress,
		IngressName: ingress.Name,
		Namespace:   ingress.Namespace,
		Annotations: ingress.Annotations,
		Result:      result,
		Log:         logger,
	}
}
