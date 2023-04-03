#!/bin/bash
if ! test -f server
then
	echo "Please compile server first (using \`go build\`)"
	echo "(this cant be done automatically as go is often installed per user)"
	exit 1
fi

if [ "$(id -u)" -ne 0 ]
then
	echo "To install please run this script as root.."
	exit 1
fi

set -e

if test -f /opt/fotos/server
then
	echo "remove old installation (executable only)"
	rm /opt/fotos/server
fi

#install into /opt
mkdir -p /opt/fotos

echo "install new executable"
cp server /opt/fotos/server
chmod +x /opt/fotos/server

if ! test -f /opt/fotos/fotos.cfg; then
	echo "install default config file"
	cp ./default.cfg /opt/fotos/fotos.cfg
fi

echo "ensuring fotos.sqlite"
chown fotos:fotos -R /opt/

echo "adding user fotos"
useradd -s /bin/nologin \
	fotos

echo "install systemd service"
cp ./fotos.service /etc/systemd/system/fotos.service

echo "reloading daemons (this can takes some time)"
systemctl daemon-reload

echo ""
echo "+TLDR---"
echo "| Installed fotos-go/server into /opt/fotos"
echo "| The configuration can be found in /opt/fotos/config"
echo "| The systemd service is called 'fotos'"
