# Install the API Gateway on Ubuntu

The following guide will walk through how to set up the API Gateway on Ubuntu
14.04.4.

## About this Guide

This guide assumes a basic familiarity with Linux and installing software on
Unix based systems.

## Install the API Gateway

The API Gateway is an OpenResty/nginx based application that manages all
requests directed at protected applications.

### Prerequisites for OpenResty

To prepare the [OpenResty](https://openresty.org) platform that the API Gateway
is built on, you will need to install a few prerequisites. The following
commands will install to build OpenResty.

```
$ sudo apt-get install libreadline-dev libncurses5-dev libpcre3-dev \
    libssl-dev perl make build-essential
```

### Download and Build OpenResty

Download OpenResty

```
$ curl -o openresty.tar.gz http://openresty.org/download/openresty-1.9.7.4.tar.gz
```

Unpack and build OpenResty

```
$ tar zxf openresty.tar.gz
$ cd openresty-1.9.7.4
$ ./configure
$ make
$ sudo make install
```

### Install LuaRocks

[LuaRocks](https://luarocks.org/) is a tool for managing libraries used in the
[Lua Programming Language](http://www.lua.org/). OpenResty is built on Lua and
the tools we will be using to interact with the MITREid Connect Server are
written in Lua. This tool will help us get future needed libraries installed.

Download LuaRocks

```
$ curl -o luarocks.tar.gz http://keplerproject.github.io/luarocks/releases/luarocks-2.3.0.tar.gz
```

Unpack and build LuaRocks

```
$ tar zxf luarocks.tar.gz
$ cd luarocks-2.3.0
$ ./configure --prefix=/usr/local/openresty/luajit \
              --with-lua=/usr/local/openresty/luajit/ \
              --lua-suffix=jit \
              --with-lua-include=/usr/local/openresty/luajit/include/luajit-2.1
$ make build
$ sudo make install
```

### Install lua-resty-openidc

[lua-resty-openidc](https://github.com/pingidentity/lua-resty-openidc) is a library for nginx implementing the OpenID Connect Relying Party (RP) and the OAuth 2.0 Resource Server (RS) functionality.

This library depends on two LuaRocks lua-resty-http and lua-resty-session. They
can be installed with the following commands:

```
$ sudo /usr/local/openresty/luajit/bin/luarocks install lua-resty-http
$ sudo /usr/local/openresty/luajit/bin/luarocks install lua-resty-session
```

Next, download and copy the lua-resty-openidc into the appropriate location.
Unfortunately, this library is not packaged as a LuaRock.

```
$ curl -o openidc.lua https://raw.githubusercontent.com/pingidentity/lua-resty-openidc/master/lib/resty/openidc.lua
$ sudo cp openidc.lua /usr/local/openresty/lualib/resty/
```
Now OpenResty/nginx have the necessary libraries to act as the API Gateway.

## Install MITREid Connect

In this setup, we will install the MITREid Connect server on the same system as
the API Gateway. This is not necessary and it is possible to install these two
components on separate systems.

### Prerequisites for MITREid Connect

MITREid Connect is a java application that uses [Maven](https://maven.apache.org/)
to manage it's build process and dependencies. The version of Maven in the
standard Ubuntu repositories is too old for MITREid Connect. Additionally, the
default version of Java is too old to run the web server component of MITREid
Connect. The following commands will install a more recent version of Java and
Maven.

```
$ sudo add-apt-repository ppa:webupd8team/java
$ sudo apt-get update
$ sudo apt-get install oracle-java8-installer
$ sudo apt-get purge maven maven2
$ sudo apt-add-repository ppa:andrei-pozolotin/maven3
$ sudo apt-get update
$ sudo apt-get install maven3
```

### Download and build MITREid Connect

You can download and build MITREid Connect by doing the following:

```
$ curl -L -o mitreid-connect.tar.gz https://github.com/mitreid-connect/OpenID-Connect-Java-Spring-Server/archive/mitreid-connect-1.2.5.tar.gz
$ tar zxf mitreid-connect.tar.gz
$ mv OpenID-Connect-Java-Spring-Server-mitreid-connect-1.2.5/ mitreid-connect
$ cd mitreid-connect
$ mvn clean install
$ cd openid-connect-server-webapp
$ mvn jetty:run
```

Note: To make things easier, I rename the folder that MITREid Connect is
unpacked into,but it is not necessary.

Once you execute the `mvn jetty:run` command, the MITREid Connect server will be
running.

## API Gateway configuration

At this point, you now have all of the software components necessary for the API
Gateway. The last thing that needs to be done is configuring nginx. This can be
done manually. It is possible to edit the file `/usr/local/openresty/nginx/conf/nginx.conf`
and include the appropriate information for the gateway. Examples on how to do
so are available [on the lua-resty-openidc site](https://github.com/pingidentity/lua-resty-openidc#sample-configuration-for-google-signin).
Note that you will have to log into the MITREid Connect server, which if you
followed the instructions above will be running on http://localhost:8080/openid-connect-server-webapp/
In the console, you can register the clients for the OIDC RP and OAuth 2 RS.
As an alternative, there is an application you can set up to automatically
generate an nginx configuration by registering clients withe MITREid Connect
using its API.

### (Optional) Run the registration application

The registration application will be cloned from a GitHub repository. It is
written in [Go](https://golang.org/) so that must be installed as well. The
following commands will set everything up:

```
$ sudo apt-get git
$ sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable
$ sudo apt-get update
$ sudo apt-get install golang
```

Go needs a place to store any downloaded libraries. Assuming you are in your
account's home directory, the following commands will set up what Go needs:

```
$ mkdir go
$ export GOPATH=~/go
```

Now, we can clone and run the registration application:

```
$ git clone https://github.com/mitre/argonaut-gateway.git
$ cd argonaut-gateway/register_nginx
$ go get
$ go run register.go
```

The `go get` command will download all of the necessary dependencies for the
registration application to run. This application must be run from session where
a browser can be launched. To access the MITREid Connect API, the registration
application itself must register as an OAuth 2 client and to do so, must
initiate a browser session to have the user approve the client.

Once the registration application has finished, it will generate an `nginx.conf`
file in the current directory. This file will have the proper configuration to
act as an OIDC RP and OAuth RS. It assumes that it is protecting an application
running on http://localhost:3001. To run nginx with this configuration, you can
do the following:

```
$ sudo cp sudo cp nginx.conf /usr/local/openresty/nginx/conf/
$ sudo /usr/local/openresty/nginx/sbin/nginx
```
