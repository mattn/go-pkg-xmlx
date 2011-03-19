# This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
# license. Its contents can be found at:
# http://creativecommons.org/publicdomain/zero/1.0/

include $(GOROOT)/src/Make.inc

TARG = github.com/jteeuwen/go-pkg-xmlx
GOFILES = document.go node.go entitymap.go

include $(GOROOT)/src/Make.pkg
