  if [ -d "organizations" ]; then
    rm -Rf organizations
  fi

  set -x
../../bin/cryptogen generate --config=./cryptogen.yaml --output=organizations
