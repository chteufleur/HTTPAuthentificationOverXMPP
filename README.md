# HTTPAuthentificationOverXMPP

Provide an HTTP anthentification over XMPP. Implementation of [XEP-0070](https://xmpp.org/extensions/xep-0070.html).


## Compilation
### Dependencies

 * [go-xmpp](https://git.kingpenguin.tk/chteufleur/go-xmpp) for the XMPP part.
 * [cfg](https://github.com/jimlawless/cfg) for the configuration file.


Download the CA at [https://kingpenguin.tk/ressources/cacert.pem](https://kingpenguin.tk/ressources/cacert.pem), then install it on your operating system.
Once installed, go into your $GOPATH directory and go get the source code.
```sh
go get git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git
```

### Configure
Configure the gateway by editing the ``httpAuth.cfg`` file in order to give all XMPP component and HTTP server informations.

XMPP
 * xmpp_server_address : Component server address connection (default: 127.0.0.1)
 * xmpp_server_port : Component server port connection (default: 5347)
 * xmpp_hostname : Component hostname
 * xmpp_secret : Component password
 * xmpp_debug : Enable debug log at true (default: false)

HTTP
 * http_port : HTTP port to bind (default: 9090)
 * http_timeoute_sec : Define a timeout if user did not give an answer to the request (default: 60)

### Utilization
To ask authorization, just send an HTTP request to the path ``/auth`` with parameters:
 * jid : JID of the user (user@host/resource or user@host)
 * domain : Domain you want to access
 * method : Method you access the domain
 * transaction_id : Transaction identifier
 * timeout : Timeout of the request in second (default : 60, max : 300)

Example:
```
GET /auth?jid=user%40host%2fresource&domain=example.org&method=POST&transaction_id=WhatEverYouWant&timeout=120 HTTP/1.1
```

This will send a request to the given JID. If the user accept, the server will return HTTP code 200, otherwise it will return HTTP code 401.

A demo version can be found at [auth.xmpp.kingpenguin.tk](http://auth.xmpp.kingpenguin.tk) for test purpose only.


## Help
To get any help, please visit the XMPP conference room at [httpauth@muc.kingpenguin.tk](xmpp:httpauth@muc.kingpenguin.tk?join) with your prefered client, or [with your browser](https://jappix.kingpenguin.tk/?r=httpauth@muc.kingpenguin.tk).
