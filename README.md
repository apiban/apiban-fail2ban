# apiban-fail2ban

**APIBAN is made possible by the generosity of our [sponsors](https://apiban.org/doc.html#sponsors).** For more information, and to get your _FREE_ APIBAN api key, please visit apiban.org.

## Contents

* [Super Simple Script Install](#super-simple-script-install)
  * [about the script](#about-the-script)
* [Get an APIBAN APIKEY](#get-an-apiban-apikey)
* [Using the Go client](#using-the-go-executable)
  * [Install Instructions](#quick-and-easy-install-instructions)
  * [Pulling all addresses](#pulling-all-addresses)
  * [Changing Data Set](#changing-data-set)
  * [Config parameters](#config-parameters)
* [Logs](#logs)
  * [Log Rotation](#log-Rotation)
* [Automation](#automation)
  * [Cron](#cron)
  * [Systemd](#systemd)
* [License / Warranty](#license--warranty)
* [Support](#support)

## Super Simple Script Install

Please at least look at the script before blindly running it on your system.

**NOTE: You need an APIKEY before running this command.**

Don't have a key? No problem. Visit [apiban.org](https://apiban.org) to get your free key.

Then, once you have your APIKEY, run:  
`curl -sSL https://raw.githubusercontent.com/apiban/apiban-fail2ban/main/install.sh | bash -s -- APIBANKEY`  
_where APIKEY is your APIBAN API KEY_

### About the script

The script will install the `apiban-fail2ban` client in `/usr/local/bin/apiban/`. The executable was compiled for amd64 architectures and will not work on pi's (you'll need to compile it yourself).

An apiban-fail2ban.service and timer are also created allowing the client to regularly check for new IP addresses. The default config created uses the jail of `asterisk-iptables` and an apiban data set of `all` (SIP and HTTP). These default values can be changed in `/usr/local/bin/apiban/config.json`.

Check out the [Using the Go client](#using-the-go-executable) section for more info on using apiban-fail2ban.

## Get an APIBAN APIKEY

Getting an APIKEY is easy and **FREE** (thanks to our sponsors).

1. Go to [apiban.org/getkey.html](https://apiban.org/getkey.html)
2. Enter your Name and Email address
3. Check your email (and spam folder) for the key.

## Using the GO executable

You can build the client using go, or just use the pre-built executable. The user running the executable will need permission to run fail2ban commands.

Be sure to update the `jail` in the config to match your desired jail.

### Quick and Easy Install Instructions

1. Create the folder `/usr/local/bin/apiban`
  
```shell 
mkdir /usr/local/bin/apiban 
```

2. Download apiban-fail2ban to `/usr/local/bin/apiban/`
    
```shell 
cd /usr/local/bin/apiban    
```

```shell 
wget https://github.com/apiban/apiban-fail2ban/raw/main/apiban-fail2ban
```

3. Download `config.json` to `/usr/local/bin/apiban/`

```shell
cd /usr/local/bin/apiban
```

```shell
wget https://github.com/apiban/apiban-fail2ban/raw/main/config.json
```

4. Using your favorite text editor, update `config.json` with your APIBAN key, for e.g:

```shell
vi config.json
```

5. Give apiban-fail2ban execute permission

```shell
chmod +x /usr/local/bin/apiban/apiban-fail2ban
```

6. Test

```shell 
/usr/local/bin/apiban/apiban-fail2ban 
```

### Pulling all addresses

Normally, apiban-fail2ban will add just the ip's that are needed to be blocked since the last successful check. Sometimes, such as after a reboot (or restart of fail2ban), you may want to pull **all** the active address. To do so, simply use the `FULL` argument. For example:

`/usr/local/bin/apiban/apiban-fail2ban FULL`

Please note, a FULL pull can take a bit to add to fail2ban.

### Changing Data Set

The default data set chosen is `all`, which incorporates both the SIP and HTTP/HTTPS honeypot data. If you wanted to have just SIP or HTTP, change the `/usr/local/bin/apiban/config.json` `set` value to either `sip`, `http`, or `all`.

### Config parameters

| parameter | description |
| --- | --- |
| `apikey` | your APIBAN [APIKEY](#get-an-apiban-apikey) |
| `lkid` | **l**ast **k**nown **id** - the "id" of the last ip address added |
| `version` | the version of the config |
| `set` | data set to use (`all`, `http`, or `sip`) |
| `flush` | used to determine when to refresh data (about 7 days from last [FULL pull](#pulling-all-addresses))
| `jail` | the fail2ban jail to add ip address |

## Logs

Log output is saved to `/var/log/apiban-client.log`. 

### Log Rotation

Want to rotate the log? Here's an example...

```bash
cat > /etc/logrotate.d/apiban-client << EOF
/var/log/apiban-client.log {
        daily
        copytruncate
        rotate 7
        compress
}
EOF
```

## Automation

### Cron

Example crontab running every 4 min...

```bash
# update apiban iptables
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin
*/4 * * * * /usr/local/bin/apiban/apiban-fail2ban >/dev/null 2>&1
```

### Systemd

Example service style automation with a 5 minute timer

```bash
cat > /lib/systemd/system/apiban-fail2ban.service << EOF
[Unit]
Description=APIBAN blocker for fail2ban
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/apiban/apiban-fail2ban

[Install]
WantedBy=multi-user.target
EOF

cat > /lib/systemd/system/apiban-fail2ban.timer << EOF
[Unit]
Description=APIBan fail2ban service schedule

[Timer]
OnUnitActiveSec=300

[Install]
WantedBy=timers.target
EOF

systemctl enable apiban-fail2ban.timer
systemctl enable apiban-fail2ban.service
systemctl start apiban-fail2ban.timer
systemctl start apiban-fail2ban.service
```

## License / Warranty

apiban-fail2ban is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 2 of the License, or (at your option) any later version

apiban-fail2ban is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

## Support

Support is provided by [LOD](https://lod.com/) and an [APIBAN room](https://matrix.to/#/#apiban:matrix.lod.com) is available on the LOD Matrix homeserver.
