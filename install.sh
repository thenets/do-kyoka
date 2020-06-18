#!/bin/bash

CONFIG_DIR=/etc/do-kyoka/
CONFIG_FILE_NAME=config.yaml
SYSTEMD_SERVICE_NAME=do-kyoka

# Download binary
rm -f /usr/local/bin/do-kyoka
cp do-kyoka /usr/local/bin/

# Create config file if not exist
mkdir -p ${CONFIG_DIR}
if [[ ! -f "${CONFIG_DIR}/${CONFIG_FILE_NAME}" ]]; then
cat <<EOT >> ${CONFIG_DIR}/${CONFIG_FILE_NAME}
apiToken: <your_digital_ocean_token_here>
firewall:
name: do-kyoka
tag: do-kyoka
ports:
  - 22
  - 80
EOT
fi

# Create systemd file
cat <<EOT > /etc/systemd/system/${SYSTEMD_SERVICE_NAME}.service
[Unit]
Description=Allows access using DigitalOcean Firewall rule
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1800
User=root
ExecStart=/usr/bin/env do-kyoka
Environment=DO_KYOKA_LOG_FILE=/var/log/do-kyoka.log
[Install]
WantedBy=multi-user.target
EOT

systemctl daemon-reload
systemctl restart ${SYSTEMD_SERVICE_NAME}
systemctl enable ${SYSTEMD_SERVICE_NAME}
systemctl reenable ${SYSTEMD_SERVICE_NAME}
