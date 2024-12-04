package clnode

/*
This file has all the defaults we set on CL node
If you see a template like {{.HTTPPort}} that means we override this dynamically inside the framework
Dynamic settings are usually used to connect Docker components together
*/

const (
	DefaultTestKeystorePassword = "thispasswordislongenough"
	DefaultPasswordTxt          = `T.tLHkcmwePT/p,]sYuntjwHKAsrhm#4eRs4LuKHwvHejWYAC2JP4M8HimwgmbaZ` //nolint:gosec
	DefaultAPICredentials       = `notreal@fakeemail.ch
fj293fbBnlQ!f9vNs`
	DefaultAPIUser     = `notreal@fakeemail.ch`
	DefaultAPIPassword = `fj293fbBnlQ!f9vNs` //nolint:gosec
)

const defaultConfigTmpl = `
[Log]
Level = 'debug'

[Pyroscope]
ServerAddress = 'http://host.docker.internal:4040'
Environment = 'local'

[WebServer]
HTTPWriteTimeout = '30s'
SecureCookies = false
HTTPPort = {{.HTTPPort}}

[WebServer.TLS]
HTTPSPort = 0

[JobPipeline]
[JobPipeline.HTTPRequest]
DefaultTimeout = '30s'
`

const dbTmpl = `[Database]
URL = '{{.DatabaseURL}}'

[Password]
Keystore = '{{.Keystore}}'
`
