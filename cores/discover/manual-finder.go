package discover

type DirectFinder struct{}

func NewDirectFinder() *DirectFinder {
	return &DirectFinder{}
}

func (df *DirectFinder) GetAddress(service string) string {
	return service
}

func (df *DirectFinder) GetAllAddress(service string) []string {
	return []string{service}
}

func (df *DirectFinder) GetAddressWithTag(service string, tag string) string {
	return df.GetAddress(service)
}

func (df *DirectFinder) GetAllAddressWithTag(service string, tag string) []string {
	return df.GetAllAddress(service)
}

func (df *DirectFinder) RegisterService(service string, address string) error {
	return nil
}

func (df *DirectFinder) RegisterServiceWithTag(service string, address string, tag string) error {
	return df.RegisterService(service, address)
}

func (df *DirectFinder) RegisterServiceWithTags(service string, address string, tags []string) error {
	return df.RegisterService(service, address)
}

func (df *DirectFinder) Close() {
	// do nothing
}
