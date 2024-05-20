DIR 		= ./build
EXECUTABLE  	= hanzi-rename
GOARCH		= amd64
GOOSWIN		= windows
GOOSX		= darwin
GOOSLINUX	= linux
GOMOD		= on
CGO_ENABLED 	= 0

WINBIN 		= $(DIR)/$(EXECUTABLE)-win-$(GOARCH).exe
OSXBIN 		= $(DIR)/$(EXECUTABLE)-darwin-$(GOARCH)
LINUXBIN 	= $(DIR)/$(EXECUTABLE)-linux-$(GOARCH)

CC 		= go build
CFLAGS		= 
LDFLAGS		= all=-w -s
GCFLAGS 	= all=
ASMFLAGS 	= all=

.PHONY: all
all: darwin linux win64

.PHONY: darwin
darwin: $(OSXBIN)
	chmod +x $(OSXBIN)

.PHONY: $(OSXBIN)
$(OSXBIN):
	GO111MODULE=$(GOMOD) GOARCH=$(GOARCH) GOOS=$(GOOSX) CGO_ENABLED=$(CGO_ENABLED) $(CC) $(CFLAGS) -o $(OSXBIN) -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)"

.PHONY: linux
linux: $(LINUXBIN)
	chmod +x $(LINUXBIN)

.PHONY: $(LINUXBIN)
$(LINUXBIN):
	GO111MODULE=$(GOMOD) GOARCH=$(GOARCH) GOOS=$(GOOSLINUX) CGO_ENABLED=$(CGO_ENABLED) $(CC) $(CFLAGS) -o $(LINUXBIN) -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)"

.PHONY: win64
win64: $(WINBIN)

.PHONY: $(WINBIN)
$(WINBIN):
	GO111MODULE=$(GOMOD) GOARCH=$(GOARCH) GOOS=$(GOOSWIN) CGO_ENABLED=$(CGO_ENABLED) $(CC) $(CFLAGS) -o $(WINBIN) -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)"

.PHONY: clean
clean:
	rm -rf $(DIR)/*
