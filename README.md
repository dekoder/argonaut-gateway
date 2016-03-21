# argonaut-gateway
Tools for building a gateway to protect healthcare applications conforming to
the [Argonauts Project](http://argonautwiki.hl7.org/).

## Overview

The goal of this project is to provide tools to create an API Gateway for
healthcare applications that wish to conform to the profiles put forward by the
Argonauts Project. Applications behind the gateway don't have implement the
the technologies for user authentication or delegated authorization. These are
handled by the gateway. If the gateway deems that a request carries the appropriate
credentials, it will pass it forward with additional HTTP headers that the
application can use.

While this approach can relieve application developers of the implementation of
things such as an OpenID Connect Relying Party infrastructure or managing
OAuth2 Bearer or Refresh tokens, it does not remove the need for the application
to make authorization decisions.

Applications will be passed information about the user or service attempting to
obtain data. It will still be the application's responsibility to determine what,
if any, data should be provided to a particular user or whether they have
permission to invoke a particular service.

This gateway simply provides applications with the information it needs to make
those decisions.

## Architecture

```
Major Components

   +----------+
   |          |
   |  Clients |
   |          |
   +-----+----+
         |
         |
         |                   +------------------------------+
+--------v---------+         |                              |
|                  |         | OpenID Connect Based         |
| Argonaut Gateway +---------+ Identity Provider            |
|                  |         | OAuth 2 Authorization Server |
+--------+---------+         |                              |
         |                   +------------------------------+
         |
  +------+------+
  |             |
  | Protected   |
  | Application |
  |             |
  +-------------+

```

In this system, the Argonaut Gateway is an instance of [nginx](http://nginx.org/)
set up with [lua-resty-openidc](https://github.com/pingidentity/lua-resty-openidc). This
allows the nginx instance to operate as an OpenID Connect Relying party.

The OpenID Connect based identity provider being used is [MITREid Connect](https://github.com/mitreid-connect).

## Protected Applications

In this architecture, applications that are protected by the API Gateway will
only receive authenticated requests. That means that the request is originating
from a user in a web browser who has logged in at the MITREid Connect acting as
an OpenID Connect IdP or it is coming from an application which has been
delegated access, with the MITREid Connect server acting as an OAuth 2
authorization server.

### Allowing users access to protected applications using OpenID Connect

Protected applications will receive proxied requests through the gateway. This
is achieved using nginx [proxy_pass](http://nginx.org/en/docs/http/ngx_http_proxy_module.html)
directive. When users make their first request to the gateway, they will be
redirected to the OpenID Connect Identity Provider to log in. Upon successful
log in, they will be redirected back to the gateway with additional information
to allow the gateway to establish identity, per the standard OpenID Connect
[Authentication using the Hybrid Flow](http://openid.net/specs/openid-connect-core-1_0.html#HybridFlowAuth).

The gateway will then obtain information about the user by accessing the
[UserInfo](http://openid.net/specs/openid-connect-core-1_0.html#UserInfo) of the
OpenID Connect Identity Provider. The claims returned by the provider will be
placed, in JSON format, into an `X-USER` header to all requests sent to the
protected application.

It is the responsibility of the protected application to determine what data and
actions should be allowed for this user based on the provided claims data.

### Allowing users access to protected applications using OAuth 2

Applications protected by the gateway for delegated access by OAuth 2 work in a
similar fashion to those allowing access via OpenID Connect. It is also possible,
and probably common, to allow access to an application using both.

For the case of delegated access using OAuth 2, we assume that a user wants to
allow a client to access data in a protected application. We also assume that
the user has an account at the OpenID Connect Identity Provider. To start the
process of a client accessing data using OAuth 2, the client should direct the
user to the OpenID Connect Identity Provider to obtain an access token. This
follows the standard OAuth 2 [Authorization Code Grant Flow](http://tools.ietf.org/html/rfc6749#section-4.1).

When the client has obtained an access token, it will make a request to the
protected application through the gateway with the access token. The gateway
will then check the validity of the access token through
[Token Introspection](https://tools.ietf.org/html/rfc7662) at the OpenID Connect
Identity Provider. If the token is valid, the Identity Provider will supply the
gateway with information about the delegated application.

Specifically, protected applications can look for two HTTP headers for requests
that have been delegated using OAuth 2:

* The `X-DELEGATED` header will be set to `true`
* The `X-SCOPE` header will contain scopes that have been grated to this client.
It is expected that these scopes will conform to [SMART on FHIR Access Scopes](http://docs.smarthealthit.org/authorization/scopes-and-launch-context/).

As with OpenID Connect, protected applications are responsible for handling the
information provided in the HTTP headers to make decisions about what actions
the client can perform.

## Installation

### Ubuntu

[Step-by-step Guide](https://github.com/mitre/argonaut-gateway/blob/master/install-ubuntu.md)

### Mac OS X

We have created [Homebrew](http://brew.sh/) formula to set up the API Gateway
and MITREid Connect servers. You can do so with the following commands, assuming
you have Homebrew installed:

```
brew tap mitre/oidc
brew install ngx_openresty
brew install openid-connect-java-spring-server
```

## License

Copyright 2016 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
