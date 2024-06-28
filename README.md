# openconnect-systemd-credentials

`openconnect-systemd-credentials` is an LD_PRELOAD library designed to supply VPN password 
from systemd-credentials store in to the OpenConnect VPN client and also optional OTP code from
Golang-based OTP helper that also uses systemd-credentials store as storage for OTPAuth URL.

## Features of LD_PRELOAD library

- Output (authentication data) of custom command in `AUTH_CMD` environment variable is fed directly into OpenConnect's stdin.

## Features of OC-OTP, Golang-based OTP helper

- CLI flag `--config` - OpenConnect VPN config file location.
- CLI flag `--otp-auth-key` - Secret key name for OTP Authentication URL.
- CLI flag `--form-entry` - OpenConnect VPN config "form-entry" key name.

## Compiling LD_PRELOAD library from source

```bash
git clone https://github.com/s3rj1k/openconnect-systemd-credentials.git
cd openconnect-systemd-credentials
make build
```

## Compiling OTP helper from source

```bash
git clone https://github.com/s3rj1k/openconnect-systemd-credentials.git
cd openconnect-systemd-credentials/go-oc-otp
make build
```

## Usage

### Systemd service without OTP

```
[Unit]
Description=OpenConnect VPN
After=network-online.target
Wants=network-online.target
ConditionUser=user

[Service]
Type=simple
KillMode=process
Restart=always
NoNewPrivileges=true

# echo "PASSWORD" | sudo systemd-creds --with-key=tpm2 encrypt --pretty --name=OC_PASSWORD - -
SetCredentialEncrypted=OC_PASSWORD: \
        DHzAexF2RZGcSwvqCLwg/iAAAAABAAAADAAAABAAAADqnHcKWcEtC4FtZJEAAAAAgAAAA \
        AAAAAALACMA8AAAACAAAAAAngAgldNq2CiPI2eIy3OAjOZlr8YCMRSogofDbVYvyR0s7v \
        0AECcGKjw4kBmpVBTPleIZoj5C8xO1b+n4BKjnUhwoGmUgFbeuCNNeh6y5D65Sir5SgIY \
        jraaBP6mT4WtUnV/D7oCttz21xzCtmfLDyqEaubfn2hxdm2Rk5CHP3qTTQ+3jKsmxOZdY \
        a0DVwHU2yKcKc9fUAVm8Sawf3PfhAE4ACAALAAAEEgAgIzKGpj9Fg1oMKJxo68mDGilVf \
        aw/kHWM8tSQi+rVB6wAEAAgycDccKN01NLFIs/Shc1JKiZIyIgIdSQWTOEX+GArKvEjMo \
        amP0WDWgwonGjryYMaKVV9rD+QdYzy1JCL6tUHrAAAAADQBUETk4CIshkQpkHsP2p4lMf \
        jnKa4PcPasVt7KgaLAX8JmRE4TP0k7Di8+LFLFLFaH5w+X6TOLik=

PrivateTmp=true
PrivateMounts=true
ProtectHome=tmpfs
ProtectSystem=strict

Environment="AUTH_CMD=/usr/bin/systemd-creds --with-key=tpm2 cat OC_PASSWORD"
Environment="LD_PRELOAD=/usr/lib/oc-password.so"

ExecStart=/usr/bin/openconnect \
        --disable-ipv6 \
        --http-auth=Basic \
        --local-hostname=user.local \
        --no-deflate \
        --no-dtls \
        --no-external-auth \
        --no-proxy \
        --non-inter \
        --os=linux-64 \
        --passwd-on-stdin \
        --protocol=anyconnect \
        --script-tun \
        --script="/usr/bin/ocproxy --dynfw 1080 --keepalive 60 --verbose" \
        --server=https://vpn.domain.net \
        --user=user \
        --verbose

[Install]
WantedBy=default.target
```

### Systemd service with OTP
```
[Unit]
Description=OpenConnect VPN
After=network-online.target
Wants=network-online.target
ConditionUser=user

[Service]
Type=simple
KillMode=process
Restart=always
NoNewPrivileges=true

# echo "PASSWORD" | sudo systemd-creds --with-key=tpm2 encrypt --pretty --name=OC_PASSWORD - -
SetCredentialEncrypted=OC_PASSWORD: \
        DHzAexF2RZGcSwvqCLwg/iAAAAABAAAADAAAABAAAADqnHcKWcEtC4FtZJEAAAAAgAAAA \
        AAAAAALACMA8AAAACAAAAAAngAgldNq2CiPI2eIy3OAjOZlr8YCMRSogofDbVYvyR0s7v \
        0AECcGKjw4kBmpVBTPleIZoj5C8xO1b+n4BKjnUhwoGmUgFbeuCNNeh6y5D65Sir5SgIY \
        jraaBP6mT4WtUnV/D7oCttz21xzCtmfLDyqEaubfn2hxdm2Rk5CHP3qTTQ+3jKsmxOZdY \
        a0DVwHU2yKcKc9fUAVm8Sawf3PfhAE4ACAALAAAEEgAgIzKGpj9Fg1oMKJxo68mDGilVf \
        aw/kHWM8tSQi+rVB6wAEAAgycDccKN01NLFIs/Shc1JKiZIyIgIdSQWTOEX+GArKvEjMo \
        amP0WDWgwonGjryYMaKVV9rD+QdYzy1JCL6tUHrAAAAADQBUETk4CIshkQpkHsP2p4lMf \
        jnKa4PcPasVt7KgaLAX8JmRE4TP0k7Di8+LFLFLFaH5w+X6TOLik=

# echo "otpauth://..." | sudo systemd-creds --with-key=tpm2 encrypt --pretty --name=OC_OTP_AUTH - -
SetCredentialEncrypted=OC_OTP_AUTH: \
        DHzAexF2RZGcSwvqCLwg/iAAAAABAAAADAAAABAAAAAKaQ2ak0C2fvEzpVsAAAAAgAAAA \
        AAAAAALACMA8AAAACAAAAAAngAgdnVTFmRp8W3aynku4wPESiKRs0ItrhrmSg6HYIhHlm \
        wAECCJaAFwm4HC5YvIxZ5AFCuURwTXAn2IuwipLK1jDTG7QIPd0K+T1azCYcLn5tEEc/k \
        o+WLXLqUSKliYAi1Niwyey/JwTc2f8WZBXcA95x86DVCT87XbPfj+S/DPjzuYLZLdSDQ0 \
        Drlk/5oRbuzFBuSu0rRJ/6wf2qaQAE4ACAALAAAEEgAgIzKGpj9Fg1oMKJxo68mDGilVf \
        aw/kHWM8tSQi+rVB6wAEAAg0gjuzGk7ebfaIzUiOyWsV6wk9CrhN2kJJqHaFxhas1AjMo \
        amP0WDWgwonGjryYMaKVV9rD+QdYzy1JCL6tUHrAAAAADDf1/+3HOWUrgpyMkoOGhehDZ \
        pPoHCg785MHkD+Rxy57ksYbntRwG7mNtMHJaaTOd5MVlWbO6lR1QWx8wWqQ==

BindPaths=/home/user/.vpn/openconnect/

PrivateTmp=true
PrivateMounts=true
ProtectHome=tmpfs
ProtectSystem=strict

Environment="AUTH_CMD=/usr/bin/systemd-creds --with-key=tpm2 cat OC_PASSWORD"
Environment="LD_PRELOAD=/usr/lib/oc-password.so"

ExecStartPre=/usr/bin/oc-otp \
        --config=/home/user/.vpn/openconnect/vpn.oc \
        --otp-auth=key:OC_OTP_AUTH \
        --form-entry=main:secondary_password

ExecStart=/usr/bin/openconnect \
        --config=/home/user/.vpn/openconnect/vpn.oc \
        --disable-ipv6 \
        --http-auth=Basic \
        --local-hostname=user.local \
        --no-deflate \
        --no-dtls \
        --no-external-auth \
        --no-proxy \
        --non-inter \
        --os=linux-64 \
        --passwd-on-stdin \
        --protocol=anyconnect \
        --script-tun \
        --script="/usr/bin/ocproxy --dynfw 1080 --keepalive 60 --verbose" \
        --server=https://vpn.domain.net \
        --user=user \
        --verbose

[Install]
WantedBy=default.target
```

### Systemd service with OTP using encrypted files (needs [Polkit rule](polkit-1/rules.d/50-systemd-credentials.rules))
```
[Unit]
Description=OpenConnect VPN
After=network-online.target
Wants=network-online.target
ConditionUser=user

[Service]
Type=simple
KillMode=process
Restart=always
NoNewPrivileges=true

BindPaths=/home/user/.vpn/openconnect/

PrivateTmp=true
PrivateMounts=true
ProtectHome=tmpfs
ProtectSystem=strict

# echo "PASSWORD" | systemd-creds --with-key=tpm2 encrypt - /home/user/.vpn/openconnect/OC_PASSWORD
Environment="AUTH_CMD=/usr/bin/systemd-creds --with-key=tpm2 decrypt /home/user/.vpn/openconnect/OC_PASSWORD -"
Environment="LD_PRELOAD=/usr/lib/oc-password.so"

# echo "otpauth://..." | systemd-creds --with-key=tpm2 encrypt - /home/user/.vpn/openconnect/OC_OTP_AUTH
ExecStartPre=/usr/bin/oc-otp \
        --config=/home/user/.vpn/openconnect/vpn.oc \
        --otp-auth=key:OC_OTP_AUTH \
        --form-entry=main:secondary_password

ExecStart=/usr/bin/openconnect \
        --config=/home/user/.vpn/openconnect/vpn.oc \
        --disable-ipv6 \
        --http-auth=Basic \
        --local-hostname=user.local \
        --no-deflate \
        --no-dtls \
        --no-external-auth \
        --no-proxy \
        --non-inter \
        --os=linux-64 \
        --passwd-on-stdin \
        --protocol=anyconnect \
        --script-tun \
        --script="/usr/bin/ocproxy --dynfw 1080 --keepalive 60 --verbose" \
        --server=https://vpn.domain.net \
        --user=user \
        --verbose

[Install]
WantedBy=default.target
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue to discuss potential improvements or features.

## License

`openconnect-systemd-credentials` is available under the MIT license. See the LICENSE file for more info.
