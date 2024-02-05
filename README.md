# OAuth Test Client

When working with OAuth2 services, especially when developing first-party
implementations, it can be helpful to have a small client application that can
test the full oauth2 flow to ensure protocol compliance and functionality.

This software should serve this purpose, providing the most minimal OAuth2
client to act as a testbed when working with OAuth implementations.

![Example success page](./.example.png)

## Support

This client exclusively supports the latest **OAuth 2.1** flow with PKCE and an
S256 challenge method. No other flows or variations of OAuth are supported.

## Usage

Copy `config.yaml.example` to `config.yaml` and fill out the endpoints and
credentials for your OAuth provider.

Then just run the program:

```
go run .
```

This will start an HTTP server on the configured port (default `8181`).
Visiting this in your browser will immediately trigger the OAuth flow.
