include .project/go-project.mk

#
# Traget environment to build.
# TARGET_ENV=prod make all
#
TARGET_ENV ?= dev

CERTS_PREFIX=test_dolly_

.PHONY: *

.SILENT:

default: help

all: clean2 gopath tools generate change_log build gen-certs test

clean2: clean
	rm -f \
		etc/dev/certs/${CERTS_PREFIX}* \
		etc/dev/certs/rootca/${CERTS_PREFIX}*

gettools:
	mkdir -p ${TOOLS_SRC}
	$(call httpsclone,${GITHUB_HOST},golang/tools,               ${TOOLS_SRC}/golang.org/x/tools,                  release-branch.go1.10)
	$(call httpsclone,${GITHUB_HOST},go-phorce/cov-report,       ${TOOLS_SRC}/github.com/go-phorce/cov-report,     master)
	$(call httpsclone,${GITHUB_HOST},go-phorce/configen,         ${TOOLS_SRC}/github.com/go-phorce/configen,       master)
	$(call httpsclone,${GITHUB_HOST},golang/lint,                ${TOOLS_SRC}/golang.org/x/lint,                   06c8688daad7faa9da5a0c2f163a3d14aac986ca)
	$(call httpsclone,${GITHUB_HOST},cloudflare/cfssl,           ${TOOLS_SRC}/github.com/cloudflare/cfssl,         ff56ab5eb62a17e335045646238665367267a678)
	$(call httpsclone,${GITHUB_HOST},mattn/goveralls,            ${TOOLS_SRC}/github.com/mattn/goveralls,          88fc0d50edb2e4cf09fe772457b17d6981826cff)
	$(call httpsclone,gopkg.in,yaml.v2,                          ${TOOLS_SRC}/gopkg.in/yaml.v2,                    31c299268d302dd0aa9a0dcf765a3d58971ac83f)
	$(call httpsclone,${GITHUB_HOST},joho/godotenv,              ${TOOLS_SRC}/github.com/joho/godotenv,            v1.2.0)
	$(call httpsclone,${GITHUB_HOST},daviddengcn/go-colortext,   ${TOOLS_SRC}/github.com/daviddengcn/go-colortext, 1.0.0)
	$(call httpsclone,${GITHUB_HOST},mattn/goreman,              ${TOOLS_SRC}/github.com/mattn/goreman,            d4c5582ffcd7d9dae49866527c7c250b04561d7e)

tools: gettools
	GOPATH=${TOOLS_PATH} go install golang.org/x/lint/golint
	GOPATH=${TOOLS_PATH} go install golang.org/x/tools/cmd/stringer
	GOPATH=${TOOLS_PATH} go install github.com/go-phorce/cov-report/cmd/cov-report
	GOPATH=${TOOLS_PATH} go install github.com/go-phorce/configen/cmd/configen
	GOPATH=${TOOLS_PATH} go install github.com/mattn/goveralls
	GOPATH=${TOOLS_PATH} go install github.com/cloudflare/cfssl/cmd/cfssl
	GOPATH=${TOOLS_PATH} go install github.com/cloudflare/cfssl/cmd/cfssljson
	GOPATH=${TOOLS_PATH} go install github.com/mattn/goreman

version:
	gofmt -r '"GIT_VERSION" -> "$(GIT_VERSION)"' version/current.template > version/current.go

build:
	echo "*** Building dolly-test"
	cd ${TEST_DIR} && go build -o ${PROJ_ROOT}/bin/dolly-test ./cmd/dolly-test

gen-certs: hsmconfig gen_test_certs
	echo "*** Running gen-certs"

hsmconfig:
	echo "*** Running hsmconfig"
	mkdir -p ~/softhsm2
	.project/config-softhsm.sh --pin-file ~/softhsm2/pin_dev.txt --generate-pin -s dolly_dev -o $(PROJ_ROOT)/etc/dev/softhsm_dev.json --list-slots --list-object
	.project/config-softhsm.sh --pin-file ~/softhsm2/pin_unittest.txt --generate-pin -s dolly_unittest -o $(PROJ_ROOT)/etc/dev/softhsm_unittest.json --delete --list-slots --list-object

gen_test_certs:
	echo "*** Running gen_test_certs in $(PROJ_ROOT) with $(CERTS_PREFIX)"
	.project/gen_test_certs.sh --hsm-cfg "$(PROJ_ROOT)/etc/dev/softhsm_dev.json" --ca-config "$(PROJ_ROOT)/etc/dev/ca-config.dev.json" --out-dir "$(PROJ_ROOT)" --prefix "$(CERTS_PREFIX)" --root-ca --ca --server --client --admin
	.project/gen_test_certs.sh --hsm-cfg "$(PROJ_ROOT)/etc/dev/softhsm_unittest.json" --ca-config "$(PROJ_ROOT)/etc/dev/ca-config.dev.json" --out-dir "$(PROJ_ROOT)" --prefix test_untrusted_ --root-ca --ca --admin

run:
	dolly_DIR=. goreman -basedir $(PROJ_DIR) -f $(PROJ_DIR)/Procfile start

change_log:
	echo "Recent changes:" > ./change_log.txt
	git log -n 20 --pretty=oneline --abbrev-commit >> ./change_log.txt
