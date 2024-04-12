# apiban-fail2ban

**APIBAN is made possible by the generosity of our [sponsors](https://apiban.org/doc.html#sponsors).**

## Super Simple Script Install

Please at least look at the script before blindly running it on your system.

**NOTE: You need an APIKEY before running this command.**

Don't have a key? No problem. Visit [apiban.org](https://apiban.org) to get your free key.

Then, once you have your APIKEY, run:  
`curl -sSL https://raw.githubusercontent.com/apiban/apiban-fail2ban/main/install.sh | bash -s -- APIBANKEY`  
_where APIKEY is your APIBAN API KEY_

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

6. Give apiban-fail2ban execute permission

```shell
chmod +x /usr/local/bin/apiban/apiban-fail2ban
```

7. Test

```shell 
./usr/local/bin/apiban/apiban-fail2ban 
```

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
