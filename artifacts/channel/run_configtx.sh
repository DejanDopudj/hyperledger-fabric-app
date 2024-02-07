 if [ -d "config" ]; then
    rm -Rf config
  fi
  mkdir config

  ../../bin/configtxgen -configPath=. -profile FourOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID testchannelid
  # ../../bin/configtxgen -configPath=. -profile FourOrgsChannel -outputCreateChannelTx ./config/channel1.tx -channelID channel1
  # ../../bin/configtxgen -configPath=. -profile FourOrgsChannel -outputCreateChannelTx ./config/channel2.tx -channelID channel2
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors-1.tx -channelID channel1 -asOrg Org1MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org2MSPanchors-1.tx -channelID channel1 -asOrg Org2MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org3MSPanchors-1.tx -channelID channel1 -asOrg Org3MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org4MSPanchors-1.tx -channelID channel1 -asOrg Org4MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors-2.tx -channelID channel2 -asOrg Org1MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org2MSPanchors-2.tx -channelID channel2 -asOrg Org2MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org3MSPanchors-2.tx -channelID channel2 -asOrg Org3MSP
  # ../../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./config/Org4MSPanchors-2.tx -channelID channel2 -asOrg Org4MSP
