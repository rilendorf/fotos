#!/bin/bash
if ! test -f server
then
	echo "Please compile server first (using \`go build\`)"
	echo "(this cant be done automatically as go is often installed per user)"
	exit 1
fi

if [ "$(id -u)" -ne 0 ]
then
	echo "To install please run this script as root."
	exit 1
fi

set -e

if test -f /opt/fotos/server
then
	echo "remove old installation (executable only)"
	rm /opt/fotos/server
fi

#install into /opt
echo "installing into /opt/fotos"
echo "ensuring /opt/fotos"
mkdir -p /opt/fotos

echo "ensuring correct ownership on /opt/fotos"
chown fotos -R /opt/fotos

echo "install new executable"
cp server /opt/fotos/server
chmod +x /opt/fotos/server

if ! test -f /opt/fotos/fotos.cfg; then
	echo "install default config file"
	cp ./default.cfg /opt/fotos/fotos.cfg
fi

USER="fotos"
if id "$USER" >/dev/null 2>&1; then
    echo "$USER user exists"
else
	echo "adding user $USER"
	/sbin/useradd -s /bin/nologin \
		"$USER"
fi

echo "install systemd service"
cp ./fotos.service /etc/systemd/system/fotos.service

echo "reloading daemons (this can takes some time)"
systemctl daemon-reload

echo ""
echo "+TLDR---"
echo "| Installed fotos/cmd/server into /opt/fotos/server"
echo "| The configuration can be found in /opt/fotos/config.cfg"
echo "| The systemd service is called 'fotos'"
