# Packaging

## Build libraries as shared

In order to packaged the project, you will need to follow those different steps.
We will build libraries as shared libraries.

The first step is to buid the GO standard library as shared.
```
$ cd $GOPATH
$ go install -buildmode=shared std
```

Next, we will build the 2 libraries to make them shared…
```
$ cd $GOPATH
$ go install -buildmode=shared -linkshared github.com/jimlawless/cfg
$ go install -buildmode=shared -linkshared git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp
```

Then finally, build the project with the option that tell to use libraries as shared and not static.
```
$ cd $GOPATH
$ go install -linkshared git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git
```

You will find shared libraries over here :
 * ``$GOROOT/pkg/linux_amd64_dynlink/libstd.so``
 * ``$GOPATH/pkg/linux_amd64_dynlink/libgithub.com-jimlawless-cfg.so``
 * ``$GOPATH/pkg/linux_amd64_dynlink/libgit.kingpenguin.tk-chteufleur-go-xmpp.git-src-xmpp.so``

And the binary file here :
 * ``$GOPATH/bin/HTTPAuthentificationOverXMPP.git``

An example can be found [here](https://github.com/jbuberel/buildmodeshared/tree/master/gofromgo).

## Configuration

At the installation, the folder ``$XDG_CONFIG_DIRS/http-auth/`` need to be created (example ``/etc/xdg/http-auth``) where to place the configuration file named ``httpAuth.conf`` and the file lang configuration named ``messages.lang``.

