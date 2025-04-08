# doh-forwarder

A simple but resilient pure-Go DNS-over-HTTPS forwarder that streams local DNS requests to any doh provider(s) supporting wireformat (defaut to quad9).    
*"Do one thing and do it well"* - we somewhat try to stick with the [suckless](https://suckless.org) philosophy.  
It can be used to bypass internet service provider's dns services, for privacy improvement, added security (quad9 configuration).

## Features
- Several DOH endpoints possible
- Single binary with no external dependencies
- Quad9 threat intelligence blocking (malware/phishing) by default
- Small & readable & easily tweakable

## Requirements
- Go (for building from source)

## Configuration
Edit `config.go` to customize before building. Default should be sane (quad9 with cloudflare as backup)


## Building & Installation
```sh
# 1. Build (produces single binary)
make

# 2. Install system-wide (default: /usr/local/bin)
sudo make install

# 3. Verify
which doh-forwarder
```

## Service Management

### For runit/sv with logging:
```bash
# 1. Create service directory
sudo mkdir -p /etc/sv/doh-forwarder/{log,env}

# 2. Create run script
sudo tee /etc/sv/doh-forwarder/run <<EOF
#!/bin/sh
exec 2>&1
exec /usr/local/bin/doh-forwarder
EOF

# 3. Create log service
sudo tee /etc/sv/doh-forwarder/log/run <<EOF
#!/bin/sh
exec svlogd -tt ./main
EOF

# 4. Set permissions and enable
sudo chmod +x /etc/sv/doh-forwarder/run /etc/sv/doh-forwarder/log/run
sudo ln -s /etc/sv/doh-forwarder /var/service/
```

### For systemd with journald logging:
```toml
# /etc/systemd/system/doh-forwarder.service
[Unit]
Description=DNS-over-HTTPS Forwarder
After=network.target

[Service]
ExecStart=/usr/local/bin/doh-forwarder
Restart=always
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```
View logs with:
```bash
journalctl -u doh-forwarder -f
```

## DNS Configuration

Preserve existing resolv.conf entries while adding local forwarder:
```bash
# Backup original

sudo cp /etc/resolv.conf /etc/resolv.conf.bak

# Add 127.0.0.1 as first nameserver
sudo sed -i '1i nameserver 127.0.0.1' /etc/resolv.conf

# Verify (should show 127.0.0.1 first)
cat /etc/resolv.conf
```

Example resulting resolv.conf:
```text
nameserver 127.0.0.1      # doh-forwarder
nameserver 192.168.1.1    # Original entries
nameserver 8.8.8.8
```

## Notes

- Some programs bypass system DNS
- NetworkManager may overwrite resolv.conf - consider:  
    ```bash
    echo "prepend domain-name-servers 127.0.0.1;" | sudo tee -a /etc/dhcp/dhclient.conf
    ```
- Test with `dig +short example.com @127.0.0.1`
- Test block list with `dig @127.0.0.1 isitblocked.org`