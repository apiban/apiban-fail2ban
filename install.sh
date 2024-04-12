#-- install script for apiban-fail2ban
echo ""
echo ""
echo " need support? https://palner.com and https://lod.com"
echo ""
echo " Copyright (C) 2024	The Palner Group, Inc. (palner.com)"
echo ""
echo " apiban-fail2ban is free software; you can redistribute it and/or modify"
echo " it under the terms of the GNU General Public License as published by"
echo " the Free Software Foundation; either version 2 of the License, or"
echo " (at your option) any later version"
echo ""
echo " apiban-fail2ban is distributed in the hope that it will be useful,"
echo " but WITHOUT ANY WARRANTY; without even the implied warranty of"
echo " MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the"
echo " GNU General Public License for more details."
echo ""
echo " You should have received a copy of the GNU General Public License"
echo " along with this program; if not, write to the Free Software"
echo " Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA"
echo ""
#-- functions
usage() {
 cat << _EOF_
Usage: ${0} APIBANKEY
...with APIBANKEY being your... apiban api key. ;)

If you need a key, please visit https://apiban.org

_EOF_
}

echo "-> checking variables"
#-- check arguments and environment
if [ "$#" -ne "1" ]; then
  echo "Expected 1 argument, got $#" >&2
  usage
  exit 2
fi
APIKEY=$1

echo "-> creating apiban directory and downloading client"
mkdir /usr/local/bin/apiban
cd /usr/local/bin/apiban
wget https://github.com/apiban/apiban-fail2ban/raw/v0.0.1/apiban-fail2ban &>/dev/null
if [ "$?" -eq "0" ]
then
  echo "  -o downloaded"
else
  echo "  -x download FAILED!!"
  exit 1
fi

echo "-> setting configuration to use your apikey"
echo "{\"apikey\":\"$APIKEY\",\"lkid\":\"100\",\"version\":\"v0.0.1\",\"set\":\"all\",\"flush\":\"200\",\"jail\":\"asterisk-iptables\"}" > config.json
chmod +x /usr/local/bin/apiban/apiban-fail2ban
echo "-> setting log rotation"
cat > /etc/logrotate.d/apiban-client << EOF
/var/log/apiban-client.log {
        daily
        copytruncate
        rotate 7
        compress
}
EOF
echo "-> setting up service"
cat /lib/systemd/system/apiban-fail2ban.service << EOF
[Unit]
Description=APIBAN blocker for fail2ban
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/apiban/apiban-fail2ban

[Install]
WantedBy=multi-user.target
EOF

cat /lib/systemd/system/apiban-fail2ban.timer << EOF
[Unit]
Description=APIBan fail2ban service schedule

[Timer]
OnUnitActiveSec=300

[Install]
WantedBy=timers.target
EOF
systemctl enable apiban-fail2ban.timer
systemctl start apiban-fail2ban.timer
systemctl start apiban-fail2ban.service
echo "-> all done."
