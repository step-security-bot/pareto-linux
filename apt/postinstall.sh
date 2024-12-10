#!/bin/bash
set -e

# Check if the script is running on Ubuntu, Debian, or Pop!_OS
if [[ -f /etc/os-release ]]; then
    . /etc/os-release
    if [[ "$ID" == "ubuntu" || "$ID" == "debian" || "$ID" == "pop" ]]; then
        # Create keyrings directory
        sudo mkdir -p --mode=0755 /usr/share/keyrings

        # Download and install GPG key
        curl -fsSL https://pkg.paretosecurity.com/paretosecurity.gpg | sudo tee /usr/share/keyrings/paretosecurity.gpg >/dev/null

        # Add Pareto repository
        echo 'deb [signed-by=/usr/share/keyrings/paretosecurity.gpg] https://pkg.paretosecurity.com/debian stable main' | sudo tee /etc/apt/sources.list.d/pareto.list

        # Check for systemd
        if command -v systemctl >/dev/null 2>&1; then
            # Create socket unit
            cat << 'EOF' | sudo tee /etc/systemd/system/pareto-linux.socket > /dev/null
[Unit]
Description=Socket for pareto-linux

[Socket]
ListenStream=/var/run/pareto-linux.sock
SocketMode=0666
Accept=no

[Install]
WantedBy=sockets.target
EOF

            # Create service unit
            cat << 'EOF' | sudo tee /etc/systemd/system/pareto-linux.service > /dev/null
[Unit]
Description=Service for pareto-linux
Requires=pareto-linux.socket

[Service]
ExecStart=/usr/bin/paretosecurity helper
User=root
Group=root
StandardInput=socket
Type=oneshot
RemainAfterExit=no

[Install]
WantedBy=multi-user.target
EOF

            # Reload systemd and enable socket
            systemctl daemon-reload
            systemctl enable pareto-linux.socket
            systemctl start pareto-linux.socket
        fi
fi

