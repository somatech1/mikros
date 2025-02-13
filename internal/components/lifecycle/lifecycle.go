package lifecycle

import (
	"context"

	"github.com/somatech1/mikros/apis/lifecycle"
	"github.com/somatech1/mikros/components/definition"
)

func OnStart(s interface{}, ctx context.Context, env definition.ServiceDeploy) error {
	// Do not execute lifecycle events in tests to force them mock features
	// that are being initialized by the service.
	if env == definition.ServiceDeploy_Test {
		return nil
	}

	if l, ok := s.(lifecycle.ServiceLifecycleStarter); ok {
		return l.OnStart(ctx)
	}

	return nil
}

func OnFinish(s interface{}, ctx context.Context, env definition.ServiceDeploy) {
	// Do not execute lifecycle events in tests to force them mock features
	// that are being initialized by the service.
	if env == definition.ServiceDeploy_Test {
		return
	}

	if l, ok := s.(lifecycle.ServiceLifecycleFinisher); ok {
		l.OnFinish(ctx)
	}
}
