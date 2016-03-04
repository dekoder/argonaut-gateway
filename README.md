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

The OpenID Connect based identity provder being used is [MITREid Connect](https://github.com/mitreid-connect).

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
