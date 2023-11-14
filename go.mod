module github.com/zc2638/swag

go 1.16

require (
	github.com/kr/pretty v0.3.0 // indirect
	github.com/modern-go/reflect2 v1.0.2
	github.com/stretchr/testify v1.8.4
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

retract v0.1.0

retract v1.2.0

// this is an incorrect release
retract v1.14.0
