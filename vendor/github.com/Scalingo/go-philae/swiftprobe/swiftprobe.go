package swiftprobe

import "github.com/ncw/swift"

type SwiftProbe struct {
	name     string
	url      string
	region   string
	tenant   string
	username string
	password string
}

func NewSwiftProbe(name, url, region, tenant, username, password string) SwiftProbe {
	return SwiftProbe{
		name:     name,
		url:      url,
		region:   region,
		tenant:   tenant,
		username: username,
		password: password,
	}
}

func (p SwiftProbe) Name() string {
	return p.name
}

func (p SwiftProbe) Check() error {
	c := swift.Connection{
		UserName: p.username,
		ApiKey:   p.password,
		AuthUrl:  p.url,
		Tenant:   p.tenant,
		Region:   p.region,
	}

	err := c.Authenticate()
	if err != nil {
		return err
	}

	return nil
}
