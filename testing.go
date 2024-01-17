package mikros

import (
	"context"

	"github.com/somatech1/mikros/components/definition"
	"github.com/somatech1/mikros/components/plugin"
	"github.com/somatech1/mikros/components/testing"
)

// ServiceTesting is an object by a Service.SetupTest call. It should be used
// to create unit tests that needs to use features (internal or external).
type ServiceTesting struct {
	svc  *Service
	test *testing.Testing
}

func setupServiceTesting(ctx context.Context, svc *Service, t *testing.Testing) *ServiceTesting {
	if svc.envs.DeploymentEnv != definition.ServiceDeploy_Test {
		return &ServiceTesting{}
	}

	svcTest := &ServiceTesting{
		svc:  svc,
		test: t,
	}

	// Sets up every plugin that needs.
	iter := svc.features.Iterator()
	for p, next := iter.Next(); next; p, next = iter.Next() {
		if featureTester, ok := p.(plugin.FeatureTester); ok {
			featureTester.Setup(ctx, t)
		}
	}

	return svcTest
}

// Teardown releases every resource allocated in the SetupTest call.
func (s *ServiceTesting) Teardown(ctx context.Context) {
	iter := s.svc.features.Iterator()
	for p, next := iter.Next(); next; p, next = iter.Next() {
		if featureTester, ok := p.(plugin.FeatureTester); ok {
			featureTester.Teardown(ctx, s.test)
		}
	}
}

// Do is a function that executes tests from inside all registered features.
func (s *ServiceTesting) Do(ctx context.Context) error {
	iter := s.svc.features.Iterator()
	for p, next := iter.Next(); next; p, next = iter.Next() {
		if featureTester, ok := p.(plugin.FeatureTester); ok {
			if err := featureTester.DoTest(ctx, s.test, s.svc.definitions.ServiceName()); err != nil {
				return err
			}
		}
	}

	return nil
}
