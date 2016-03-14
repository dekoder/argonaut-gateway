This directory contains a Go application that will obtain permission to access
FHIR resources using [SMART Authorization](http://docs.smarthealthit.org/authorization/).
It is purely a demonstration application.

This application makes a few assumptions:
* The [MITREid Connect Server](https://github.com/mitreid-connect/OpenID-Connect-Java-Spring-Server) is running on localhost:8080
* The nginx gateway is running on localhost:5000
* This client has been registered with the MITREid Connect Server (you will need
  to modify the source code)

This application simply takes the necessary action to obtain a Bearer Token and
then passes through the FHIR response it receives. The application runs on
localhost:3000.
