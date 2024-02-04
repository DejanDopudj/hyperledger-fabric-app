 if [ -d "config" ]; then
    rm -Rf config
  fi
  mkdir config

  ../../bin/configtxgen -profile FourOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID testchannelid
