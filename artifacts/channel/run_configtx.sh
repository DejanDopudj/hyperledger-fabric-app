 if [ -d "config" ]; then
    rm -Rf config
  fi
  mkdir config

  ../../bin/configtxgen -profile FourOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID testchannelid
  ../../bin/configtxgen -profile FourOrgsChannel1 -outputCreateChannelTx ./config/channel1.tx -channelID channel1
  ../../bin/configtxgen -profile FourOrgsChannel2 -outputCreateChannelTx ./config/channel2.tx -channelID channel2
  ../../bin/configtxgen -profile FourOrgsChannel1 -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID channel1 -asOrg Org1MSP
  ../../bin/configtxgen -profile FourOrgsChannel1 -outputAnchorPeersUpdate ./config/Org2MSPanchors.tx -channelID channel1 -asOrg Org2MSP
  ../../bin/configtxgen -profile FourOrgsChannel1 -outputAnchorPeersUpdate ./config/Org3MSPanchors.tx -channelID channel1 -asOrg Org3MSP
  ../../bin/configtxgen -profile FourOrgsChannel1 -outputAnchorPeersUpdate ./config/Org4MSPanchors.tx -channelID channel1 -asOrg Org4MSP
  ../../bin/configtxgen -profile FourOrgsChannel2 -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID channel2 -asOrg Org1MSP
  ../../bin/configtxgen -profile FourOrgsChannel2 -outputAnchorPeersUpdate ./config/Org2MSPanchors.tx -channelID channel2 -asOrg Org2MSP
  ../../bin/configtxgen -profile FourOrgsChannel2 -outputAnchorPeersUpdate ./config/Org3MSPanchors.tx -channelID channel2 -asOrg Org3MSP
  ../../bin/configtxgen -profile FourOrgsChannel2 -outputAnchorPeersUpdate ./config/Org4MSPanchors.tx -channelID channel2 -asOrg Org4MSP
