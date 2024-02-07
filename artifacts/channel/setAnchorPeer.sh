#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# import utils
. ./envVar.sh

fetchChannelConfig() {
  ORG=$1
  CHANNEL=$2
  OUTPUT=$3

  infoln "Fetching the most recent configuration block for the channel"
  set -x
  ../../bin/peer channel fetch config config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL --tls --cafile $ORDERER_CA
  { set +x; } 2>/dev/null

  infoln "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  ../../bin/configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  { set +x; } 2>/dev/null
}

# createConfigUpdate <channel_id> <original_config.json> <modified_config.json> <output.pb>
# Takes an original and modified config, and produces the config update tx
# which transitions between the two
# NOTE: this must be run in a CLI container since it requires ../../bin/configtxlator 
createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4

  set -x
  ../../bin/configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb
  ../../bin/configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb
  ../../bin/configtxlator compute_update --channel_id "${CHANNEL}" --original original_config.pb --updated modified_config.pb >config_update.pb
  ../../bin/configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json
  ../../bin/configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}"
  { set +x; } 2>/dev/null
}

# signConfigtxAs../../bin/peerOrg <org> <configtx.pb>
# Set the ../../bin/peerOrg admin of an org and sign the config update
signConfigtxAs../../bin/peerOrg() {
  ORG=$1
  CONFIGTXFILE=$2
  setGlobals $ORG
  set -x
  ../../bin/peer channel signconfigtx -f "${CONFIGTXFILE}"
  { set +x; } 2>/dev/null
}

# NOTE: this must be run in a CLI container since it requires jq and configtxlator 
createAnchorPeerUpdate() {
  infoln "Fetching channel config for channel $CHANNEL_NAME"
  fetchChannelConfig $ORG $CHANNEL_NAME Org${ORG}MSPconfig.json

  infoln "Generating anchor peer update transaction for Org${ORG} on channel $CHANNEL_NAME"

  if [ $ORG -eq 1 ]; then
    HOST="peer0.org1.example.com"
    PORT=7051
  elif [ $ORG -eq 2 ]; then
    HOST="peer0.org2.example.com"
    PORT=7051
  elif [ $ORG -eq 3 ]; then
    HOST="peer0.org2.example.com"
    PORT=7051
  elif [ $ORG -eq 4 ]; then
    HOST="peer0.org3.example.com"
    PORT=7051
  else
    errorln "Org${ORG} unknown"
  fi

  set -x
  # Modify the configuration to append the anchor peer 
  jq '.channel_group.groups.Application.groups.Org'${ORG}'MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "'$HOST'","port": '$PORT'}]},"version": "0"}}' Org${ORG}MSPconfig.json > Org${ORG}MSPmodified_config.json
  { set +x; } 2>/dev/null

  # Compute a config update, based on the differences between 
  # {orgmsp}config.json and {orgmsp}modified_config.json, write
  # it as a transaction to {orgmsp}anchors.tx
  createConfigUpdate ${CHANNEL_NAME} Org${ORG}MSPconfig.json Org${ORG}MSPmodified_config.json Org${ORG}MSPanchors.tx
}

updateAnchorPeer() {
  ../../bin/peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME -f Org${ORG}MSPanchors.tx --tls --cafile ./organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem >&log.txt
  res=$?
  cat log.txt
  verifyResult $res "Anchor peer update failed"
  successln "Anchor peer set for org '$CORE_PEER_LOCALMSPID' on channel '$CHANNEL_NAME'"
}

ORG=$1
CHANNEL_NAME=$2
FABRIC_CFG_PATH=.
setGlobalsCLI $ORG

createAnchorPeerUpdate 

updateAnchorPeer 