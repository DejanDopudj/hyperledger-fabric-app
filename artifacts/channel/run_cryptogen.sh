  if [ -d "crypto-config" ]; then
    rm -Rf crypto-config
  fi

  set -x
../../bin/cryptogen generate --config=./cryptogen.yaml
