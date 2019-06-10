#!/bin/bash

#
# gen_certs.sh
#   --hsm-cfg {file}    - hsm config
#   --out-dir {dir}     - specify output folder
#   --prefix {prefix}   - specify prefix for files, by default: ${PREFIX}
#   --ca-config {file}  -
#   --root-ca
#   --ca
#   --server
#   --admin
#   --client
#   --peers
#   --encrypt
#
POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -o|--out-dir)
    OUT_DIR="$2"
    shift # past argument
    shift # past value
    ;;
    -h|--hsm-cfg)
    HSM_CFG="$2"
    shift # past argument
    shift # past value
    ;;
    -p|--prefix)
    PREFIX="$2"
    shift # past argument
    shift # past value
    ;;
    -c|--ca-config)
    CA_CONFIG="$2"
    shift # past argument
    shift # past value
    ;;
    --root-ca)
    ROOTCA=YES
    shift # past argument
    ;;
    --ca)
    CA=YES
    shift # past argument
    ;;
    --server)
    SERVER=YES
    shift # past argument
    ;;
    --admin)
    ADMIN=YES
    shift # past argument
    ;;
    --client)
    CLIENT=YES
    shift # past argument
    ;;
    --peers)
    PEERS=YES
    shift # past argument
    ;;
    --encrypt)
    ENCRYPT=YES
    shift # past argument
    ;;
    *)
    echo "Unsupported flag: $key"
    exit 1
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

[ -z "$OUT_DIR" ] &&  echo "Specify --out-dir" && exit 1
[ -z "$HSM_CFG" ] &&  echo "Specify --hsm-cfg" && exit 1

[ -z "$CA_CONFIG" ] && CA_CONFIG=$(OUT_DIR)/etc/dev/ca-config.dev.json
[ -z "$PREFIX" ] && PREFIX=test_

HOSTNAME=`hostname`

echo OUT_DIR   = "${OUT_DIR}"
echo CA_CONFIG = "${CA_CONFIG}"
echo PREFIX    = "${PREFIX}"

if [[ "$ROOTCA" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/rootca/${PREFIX}root_CA-key.pem ]]; then
	echo "*** generating ${PREFIX}root_CA"
    mkdir -p ${OUT_DIR}/etc/dev/certs/rootca
	cfssl genkey -initca -config=${CA_CONFIG} ${OUT_DIR}/etc/dev/certs/csr/${PREFIX}root_CA.json | cfssljson -bare ${OUT_DIR}/etc/dev/certs/rootca/${PREFIX}root_CA
fi

if [[ "$CA" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA-key.pem ]]; then
	echo "*** generating CA cert"
    mkdir -p ${OUT_DIR}/etc/dev/certs
  	cfssl genkey -initca -config=${CA_CONFIG} ${OUT_DIR}/etc/dev/certs/csr/${PREFIX}issuer_CA.json | cfssljson -bare ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA

	cfssl sign \
        -config=${CA_CONFIG} \
        -profile=CA \
        -ca ${OUT_DIR}/etc/dev/certs/rootca/${PREFIX}root_CA.pem \
        -ca-key ${OUT_DIR}/etc/dev/certs/rootca/${PREFIX}root_CA-key.pem \
        -csr ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.csr | cfssljson -bare ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA
fi

if [[ "$SERVER" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/${PREFIX}server-key.pem ]]; then
	echo "*** generating server cert"
    mkdir -p ${OUT_DIR}/etc/dev/certs
	cfssl gencert \
        -config=${CA_CONFIG} \
        -profile=server \
        -ca ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem \
        -ca-key ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA-key.pem \
        -hostname=localhost,127.0.0.1,10.77.77.100,10.77.77.101,10.77.77.102,10.88.88.100,10.88.88.101,10.88.88.102,*.dolly.in,${HOSTNAME} \
        ${OUT_DIR}/etc/dev/certs/csr/${PREFIX}server.json | cfssljson -bare ${OUT_DIR}/etc/dev/certs/${PREFIX}server
        cat ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem >> ${OUT_DIR}/etc/dev/certs/${PREFIX}server.pem
fi

if [[ "$ADMIN" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/${PREFIX}admin-key.pem ]]; then
	echo "*** generating admin cert"
    mkdir -p ${OUT_DIR}/etc/dev/certs
	cfssl gencert \
        -config=${CA_CONFIG} \
        -profile=client \
        -ca ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem \
        -ca-key ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA-key.pem \
        ${OUT_DIR}/etc/dev/certs/csr/${PREFIX}admin.json | cfssljson -bare ${OUT_DIR}/etc/dev/certs/${PREFIX}admin
        cat ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem >> ${OUT_DIR}/etc/dev/certs/${PREFIX}admin.pem
fi

if [[ "$CLIENT" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/${PREFIX}client-key.pem ]]; then
	echo "*** generating client cert"
    mkdir -p ${OUT_DIR}/etc/dev/certs
	cfssl gencert \
        -config=${CA_CONFIG} \
        -profile=client \
        -ca ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem \
        -ca-key ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA-key.pem \
        ${OUT_DIR}/etc/dev/certs/csr/${PREFIX}client.json | cfssljson -bare ${OUT_DIR}/etc/dev/certs/${PREFIX}client
        cat ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem >> ${OUT_DIR}/etc/dev/certs/${PREFIX}client.pem
fi

if [[ "$PEERS" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/${PREFIX}peers-key.pem ]]; then
	echo "*** generating peers cert"
    mkdir -p ${OUT_DIR}/etc/dev/certs
	cfssl gencert \
        -config=${CA_CONFIG} \
        -profile=peers \
        -ca ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem \
        -ca-key ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA-key.pem \
        -hostname=localhost,127.0.0.1,10.77.77.100,10.77.77.101,10.77.77.102,10.88.88.100,10.88.88.101,10.88.88.102,*.dolly.in,${HOSTNAME} \
        ${OUT_DIR}/etc/dev/certs/csr/${PREFIX}peers.json | cfssljson -bare ${OUT_DIR}/etc/dev/certs/${PREFIX}peers
        cat ${OUT_DIR}/etc/dev/certs/${PREFIX}issuer_CA.pem >> ${OUT_DIR}/etc/dev/certs/${PREFIX}peers.pem
fi

if [[ "$ENCRYPT" == "YES" && ! -f ${OUT_DIR}/etc/dev/certs/${PREFIX}encrypt-key.pem ]]; then
	echo "*** generating encryption key"
    mkdir -p ${OUT_DIR}/etc/dev/certs
	enrollmepki hsm genkey \
      --hsm-cfg ${HSM_CFG} \
      --label ${PREFIX}encrypt-* \
      --purpose encryption \
      --alg RSA \
      --size 4096 \
      --output ${OUT_DIR}/etc/dev/certs/${PREFIX}encrypt-key.pem
fi
