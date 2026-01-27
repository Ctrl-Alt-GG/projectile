#!/bin/bash

mkdir /tmp/projectile
pushd /tmp/projectile

htpasswd -nb "secret-id" "secret-key" > .htpasswd

# https://letsencrypt.org/docs/certificates-for-localhost/
openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")


cat > /tmp/projectile/config.yaml << EOF
gameData:
  game: "dummy"
  name: "Dummy Game"
server:
  address: "localhost:50051"
  id: "secret-id"
  key: "secret-key"
scraper:
  module: "dummy"
EOF