FROM debian:stretch
RUN mkdir -p /etc/lora-mqtt
VOLUME /etc/lora-mqtt
WORKDIR /root/
COPY lora-mqtt-linux-arm .
ENTRYPOINT ["./lora-mqtt-linux-arm", "--config", "/etc/lora-mqtt/lora-mqtt.yaml" ,"-v"]
