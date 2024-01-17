package plugin

type ServiceSet struct {
	services map[string]Service
}

func NewServiceSet() *ServiceSet {
	return &ServiceSet{
		services: make(map[string]Service),
	}
}

func (s *ServiceSet) Register(svc Service) {
	if name := svc.Name(); name != "" {
		if _, ok := s.services[name]; !ok {
			s.services[name] = svc
		}
	}
}

func (s *ServiceSet) Services() map[string]Service {
	svc := make(map[string]Service)
	for k, v := range s.services {
		svc[k] = v
	}

	return svc
}

func (s *ServiceSet) Append(services *ServiceSet) {
	for k, v := range services.Services() {
		if _, ok := s.services[k]; !ok {
			s.services[k] = v
		}
	}
}
