#
# Copyright (c) 2018 Vivotek Inc. All rights reserved.
#
# +-----------------------------------------------------------------+
# | THIS SOFTWARE IS FURNISHED UNDER A LICENSE AND MAY ONLY BE USED |
# | AND COPIED IN ACCORDANCE WITH THE TERMS AND CONDITIONS OF SUCH  |
# | A LICENSE AND WITH THE INCLUSION OF THE THIS COPY RIGHT NOTICE. |
# | THIS SOFTWARE OR ANY OTHER COPIES OF THIS SOFTWARE MAY NOT BE   |
# | PROVIDED OR OTHERWISE MADE AVAILABLE TO ANY OTHER PERSON. THE   |
# | OWNERSHIP AND TITLE OF THIS SOFTWARE IS NOT TRANSFERRED.        |
# |                                                                 |
# | THE INFORMATION IN THIS SOFTWARE IS SUBJECT TO CHANGE WITHOUT   |
# | ANY PRIOR NOTICE AND SHOULD NOT BE CONSTRUED AS A COMMITMENT BY |
# | VIVOTEK INC.                                                    |
# +-----------------------------------------------------------------+
#
# Date: 2018-05-29 15:10:0
# Author: Willy Ting (willy.ting@vivotek.com)
#
PKG := gachamachine

PKGDIR := $(GOPATH)/src/$(PKG)
OUTPUTDIR := $(GOPATH)/src/$(PKG)
BINARY := $(OUTPUTDIR)/$(PKG)
export OUTPUTDIR

HOST_NCPU ?= 1
export HOST_NCPU

PHONY += all
all: dep
	go build -v -o $(BINARY) $(PKG)

PHONY += clean
clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

PHONY += distclean
distclean: clean
	go clean

PHONY += romfs
romfs: all

dep:
	cd $(PKGDIR)
	$(GOPATH)/bin/dep ensure
	cd -

.PHONY: $(PHONY)
