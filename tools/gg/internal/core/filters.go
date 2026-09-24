package core

func (Core) InstanceFilter() string {
	return `service=~"$service", pod=~"$pod"`
}

func (Core) LogInstanceFilter() string {
	return `service: ($service) pod: ($pod)`
}

func (Core) TraceInstanceFilter() string {
	return `{resource.service.name=$service}`
}
