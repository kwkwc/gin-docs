package gin_docs

type Config struct {
	// Title, default `API Doc`
	Title string
	// Version, default `1.0.0`
	Version string
	// Description
	Description string

	// Custom CDN CSS Template
	CdnCssTemplate string
	// Custom CDN JS Template
	CdnJsTemplate string

	// Custom url prefix, default `/docs/api`
	UrlPrefix string
	// No document text, default `No documentation found for this API`
	NoDocText string
	// Enable document pages, default `true`
	Enable bool
	// Using CDN, default `false`
	Cdn bool
	// API package name to exclude
	Exclude []string
	// Methods allowed to be displayed, default `[]string{"GET", "POST", "PUT", "DELETE", "PATCH"}`
	MethodsList []string
	// SHA256 encrypted authorization password, e.g. here is admin
	// echo -n admin | shasum -a 256
	// `8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918`
	PasswordSha2 string
	// Enable markdown processing for all documents, default `true`
	AllMd bool
}

func (c *Config) Default() *Config {
	c.Title = "API Doc"
	c.Version = "1.0.0"
	c.UrlPrefix = "/docs/api"
	c.NoDocText = "No documentation found for this API"
	c.MethodsList = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	c.Enable = true
	c.AllMd = true

	return c
}
