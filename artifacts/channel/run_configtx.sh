 if [ -d "config" ]; then
    rm -Rf config
  fi
  mkdir config

  ../../bin/configtxgen -profile FourOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID testchannelid
  #../../bin/configtxgen -profile FourOrgsOrdererGenesis -outputCreateChannelTx ./config/channel.tx -channelID $CHANNEL_NAME
  #../../bin/configtxgen -profile FourOrgsOrdererGenesis -outputAnchorPeersUpdate ./config/${MSP_NAME}anchors.tx -channelID $CHANNEL_NAME -asOrg $MSP_NAME
