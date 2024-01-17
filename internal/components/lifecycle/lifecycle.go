package lifecycle

import (
	"context"

	"github.com/somatech1/mikros/apis/lifecycle"
)

func OnStart(s interface{}, ctx context.Context) error {
	if l, ok := s.(lifecycle.ServiceLifecycleStarter); ok {
		return l.OnStart(ctx)
	}

	return nil
}

func OnFinish(s interface{}, ctx context.Context) {
	if l, ok := s.(lifecycle.ServiceLifecycleFinisher); ok {
		l.OnFinish(ctx)
	}
}
