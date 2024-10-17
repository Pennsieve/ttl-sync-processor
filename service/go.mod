module github.com/pennsieve/ttl-sync-processor/service

go 1.21

replace github.com/pennsieve/ttl-sync-processor/client => ./../client

require (
	github.com/google/uuid v1.6.0
	github.com/pennsieve/processor-post-metadata/client v0.0.0-20241017161214-7ed76902fd93
	github.com/pennsieve/processor-pre-metadata/client v0.0.0-20241015190912-3e004748369b
	github.com/pennsieve/ttl-sync-processor/client v0.0.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
