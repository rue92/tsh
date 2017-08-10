# Twitch Shell (TSH) #

## Overview ##
The intention of Twitch Shell (TSH) is to provide a standalone binary
for basic interactions with the Twitch streaming community. Though
called a shell, TSH is really more of a fancy CLI application built on
top of [termui](https://github.com/gizak/termui) (which is a fantastic
little library, by the way). This is still largely a toy project with
no concrete goals though the roadmap below does provide intentions for
the future. At some point I may move the twitch package out of this
repo into its own as it seems to be the first Go implementation for
the V5 API -- I could be wrong though! That would only happen as the 
package reaches a more mature and feature-complete stage however.

## Roadmap ##
  * OAuth authentication in order to retrieve channels and games the
    user is following
  * Search capability for channels and games
  * Support for launching Streamlink from TSH
  * Use [glide](https://github.com/Masterminds/glide) to manage
    package versions (which for now is just
    [termui](https://github.com/gizak/termui) though)
  * ??? Probably more asynchronous communication to make everything
    more responsive
    
## Building ##
Depending on how you intend to run this it can be as easy as `cd`'ing
into `$GOPATH/src/github.com/rue92/tsh` and doing
```
go run tsh.go user_config.go
```
or alternatively 
```
go build && ./tsh
```

### Installing ###
If for some reason you want to install this to `$GOPATH/bin` then go
ahead and do `go install`

