#!/bin/bash

. ./envVar.sh

CHANNEL_NAME="$1"
DELAY="$2"
MAX_RETRY="$3"
VERBOSE="$4"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}

createChannelTx() {
	set -x
	../../bin/configtxgen -profile FourOrgsChannel -configPath . -outputCreateChannelTx ./config/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME 
	res=$?
	{ set +x; } 2>/dev/null
  #verifyResult $res "Failed to generate channel configuration transaction..."
}

createChannel() {
	setGlobals 1
	local rc=1
	local COUNTER=1
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
		sleep $DELAY
		set -x
		../../bin/peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.example.com -f ./config/${CHANNEL_NAME}.tx --outputBlock $BLOCKFILE --tls --cafile ./organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem >&log.txt
		res=$?
		{ set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	#verifyResult $res "Channel creation failed"
}

# joinChannel ORG
joinChannel() {
  ORG=$1
  setGlobals $ORG
	local rc=1
	local COUNTER=1
	## Sometimes Join takes time, hence retry
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    ../../bin/peer channel join -b $BLOCKFILE >&log2.txt
    res=$?
    { set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log2.txt
	#verifyResult $res "After $MAX_RETRY attempts, peer0.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}

setAnchorPeer() {
  ORG=$1
  bash ./setAnchorPeer.sh $ORG $CHANNEL_NAME 
}

FABRIC_CFG_PATH=./channel
CHANNEL_NAME=channel1
createChannelTx
FABRIC_CFG_PATH=.
BLOCKFILE="./config/${CHANNEL_NAME}.block"
createChannel
joinChannel 1
successln "Channel '$CHANNEL_NAME' created"

FABRIC_CFG_PATH=./channel
CHANNEL_NAME=channel2
createChannelTx
FABRIC_CFG_PATH=.
BLOCKFILE="./config/${CHANNEL_NAME}.block"
createChannel
joinChannel 1
successln "Channel '$CHANNEL_NAME' created"


# ## Set the anchor peers for each org in the channel
# infoln "Setting anchor peer for org1..."
# setAnchorPeer 1
# infoln "Setting anchor peer for org2..."
# setAnchorPeer 2

#successln "Channel '$CHANNEL_NAME' joined"