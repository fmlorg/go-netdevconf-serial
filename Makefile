PWD_DIR   != pwd -P
PKG_DIR    = ${PWD_DIR}/cmd/netdevconf
BIN_DIR    = ${PWD_DIR}/bin
PROG_LINUX = ${BIN_DIR}/netdevconf
PROG_WIN11 = ${BIN_DIR}/netdevconf.exe
GO         = go

all: go

go: build
	@ test -x ${PROG_LINUX} && sudo -E ${PROG_LINUX}

build: build-linux

build-linux:
	@ ${GO} build -o ${PROG_LINUX} ${PKG_DIR}
	@ ls -l ${PROG_LINUX}

test:
	@ ${GO} test -v

console:
	@ sudo minicom -D /dev/ttyUSB0

clean:
	@ rm -vfr bin


#
# windos port
#
win:             win11
win11:           build-windows11
build-win11:     build-windows11
build-windows11:
	@ env GOOS=windows GOARCH=amd64 ${GO} build -o ${PROG_WIN11} ${PKG_DIR}
