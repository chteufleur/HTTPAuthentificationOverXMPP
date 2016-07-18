# HTTPAuthentificationOverXMPP

Provide an HTTP anthentification over XMPP. Implementation of [XEP-0070](https://xmpp.org/extensions/xep-0070.html).


## Compilation
### Dependencies

 * [go-xmpp](https://git.kingpenguin.tk/chteufleur/go-xmpp) for the XMPP part.
 * [cfg](https://github.com/jimlawless/cfg) for the configuration file.

You must first [install go environment](https://golang.org/doc/install) on your system.
Then, go into your $GOPATH directory and go get the source code.
```sh
go get git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git
```

### Configure
Configure the gateway by editing the ``httpAuth.cfg`` file in order to give all XMPP component and HTTP server informations.
An example of the config file can be found in [the repos](https://git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP/src/master/httpAuth.cfg).

XMPP
 * xmpp_server_address : Component server address connection (default: 127.0.0.1)
 * xmpp_server_port : Component server port connection (default: 5347)
 * xmpp_hostname : Component hostname
 * xmpp_secret : Component password
 * xmpp_debug : Enable debug log at true (default: false)

HTTP
 * http_port : HTTP port to bind (default: 9090, desactive: -1)
 * https_port : HTTPS port to bind (default: 9093, desactive: -1)
 * https_cert_path : Path to the certificate file (default: ./cert.pem)
 * https_key_path : Path to the key file (default: ./key.pem)
 * http_timeoute_sec : Define a timeout if user did not give an answer to the request (default: 60)


### Utilization
To ask authorization, just send an HTTP request to the path ``/auth`` with parameters:
 * __jid__ : JID of the user (user@host/resource or user@host)
 * __domain__ : Domain you want to access
 * __method__ : Method you access the domain
 * transaction_id : Transaction identifier (auto generated if not provide)
 * timeout : Timeout of the request in second (default : 60, max : 300)

__Bold parameters__ are mandatory.

Example:
```
GET /auth?jid=user%40host%2fresource&domain=example.org&method=POST&transaction_id=WhatEverYouWant&timeout=120 HTTP/1.1
```

This will send a request to the given JID, then return HTTP code depending on what appended.
 * 200 : User accept the request
 * 400 : One or more mandatory parameter(s) is missing
 * 401 : User deny the request or timeout
 * 520 : Unknown error append
 * 523 : Server is unreachable


If the provided JID contain a resource, it will try to send an ``iq`` stanza.
If the answer to this ``iq`` is a ``feature-not-implemented`` or ``service-unavailable`` error,
it will automatically send a ``message`` stanza. Unfortunately, if a ``message`` stanza is used,
their is probably no way to get the error if the JID does not exist or is unreachable.


A demo version can be found at [auth.xmpp.kingpenguin.tk](http://auth.xmpp.kingpenguin.tk) for test purpose only.


## Help
To get any help, please visit the XMPP conference room at [httpauth@muc.kingpenguin.tk](xmpp://httpauth@muc.kingpenguin.tk?join) with your prefered client, or [with your browser](https://jappix.kingpenguin.tk/?r=httpauth@muc.kingpenguin.tk).
