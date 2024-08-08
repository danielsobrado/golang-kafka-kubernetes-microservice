#!/bin/bash

set -e

# Create the kafka-creds directory if it doesn't exist
mkdir -p kafka-creds
cd kafka-creds

# Remove all existing files in the kafka-creds directory
echo "Removing existing credentials..."
rm -f *.crt *.csr *.jks *.key *.p12 *.pem *.srl *_creds client.properties

# Generate CA key and certificate
echo "Generating CA key and certificate..."
openssl req -new -nodes \
    -x509 \
    -days 365 \
    -newkey rsa:2048 \
    -keyout ca.key \
    -out ca.crt \
    -subj "/CN=ca.example.com/OU=TEST/O=EXAMPLE/L=LONDON/ST=LONDON/C=UK"

# Convert CA files to PEM format
cat ca.crt ca.key > ca.pem

# Function to create credentials for a broker
create_broker_creds() {
    BROKER=$1
    
    echo "Creating credentials for ${BROKER}..."
    
    # Create server key and certificate signing request (CSR)
    openssl req -new \
        -newkey rsa:2048 \
        -keyout ${BROKER}.key \
        -out ${BROKER}.csr \
        -subj "/CN=${BROKER}/OU=TEST/O=EXAMPLE/L=LONDON/ST=LONDON/C=UK" \
        -nodes

    # Create a configuration file for the SAN extension
    cat > ${BROKER}.ext << EOF
subjectAltName = DNS:${BROKER},DNS:localhost,DNS:kafka
EOF

    # Sign the certificate with the CA, including the SAN extension
    openssl x509 -req \
        -days 3650 \
        -in ${BROKER}.csr \
        -CA ca.crt \
        -CAkey ca.key \
        -CAcreateserial \
        -out ${BROKER}.crt \
        -extfile ${BROKER}.ext

    # Convert the server certificate to PKCS12 format
    openssl pkcs12 -export \
        -in ${BROKER}.crt \
        -inkey ${BROKER}.key \
        -chain \
        -CAfile ca.pem \
        -name ${BROKER} \
        -out ${BROKER}.p12 \
        -password pass:confluent
    
    # Create the broker keystore and import the certificate
    keytool -importkeystore \
        -deststorepass confluent \
        -destkeystore kafka.${BROKER}.keystore.pkcs12 \
        -srckeystore ${BROKER}.p12 \
        -deststoretype PKCS12 \
        -srcstoretype PKCS12 \
        -noprompt \
        -srcstorepass confluent
    
    # Create credentials files
    echo "confluent" > kafka_${BROKER}_sslkey_creds
    echo "confluent" > kafka_${BROKER}_keystore_creds

    # Clean up temporary files
    rm ${BROKER}.ext
}

# Create credentials for kafka-1, kafka-2, and kafka-3
create_broker_creds "kafka-1"
create_broker_creds "kafka-2"
create_broker_creds "kafka-3"

# Function to create or update truststore
create_or_update_truststore() {
    TRUSTSTORE=$1
    echo "Creating/updating truststore: ${TRUSTSTORE}"
    if [ -f $TRUSTSTORE ]; then
        keytool -delete -alias CARoot -keystore $TRUSTSTORE -storepass confluent -noprompt || true
    fi
    keytool -keystore $TRUSTSTORE -alias CARoot -import -file ca.crt -storepass confluent -keypass confluent -noprompt
}

# Create or update the client and server truststores
create_or_update_truststore "kafka.client.truststore.jks"
create_or_update_truststore "kafka.server.truststore.jks"

# Create truststore credentials file
echo "confluent" > kafka_truststore_creds

# Create client.properties file
cat << EOF > client.properties
security.protocol=SSL
ssl.truststore.location=/etc/kafka/secrets/kafka.client.truststore.jks
ssl.truststore.password=confluent
ssl.keystore.location=/etc/kafka/secrets/kafka.client.keystore.jks
ssl.keystore.password=confluent
ssl.key.password=confluent
EOF

# Create client keystore (JKS format)
echo "Creating client keystore..."
keytool -genkey -keyalg RSA -alias localhost \
    -keystore kafka.client.keystore.jks \
    -storepass confluent \
    -keypass confluent \
    -validity 365 \
    -keysize 2048 \
    -dname "CN=kafka-client,OU=TEST,O=EXAMPLE,L=LONDON,ST=LONDON,C=UK" \
    -storetype JKS

# Generate client certificate signing request
keytool -certreq -alias localhost -file client.csr \
    -keystore kafka.client.keystore.jks \
    -storepass confluent \
    -keypass confluent

# Sign the client certificate
openssl x509 -req -CA ca.crt -CAkey ca.key -in client.csr \
    -out client-signed.crt -days 3650 -CAcreateserial

# Import the CA cert and the signed client cert back into the client keystore
keytool -keystore kafka.client.keystore.jks -alias CARoot \
    -import -file ca.crt -storepass confluent -keypass confluent -noprompt
keytool -keystore kafka.client.keystore.jks -alias localhost \
    -import -file client-signed.crt -storepass confluent -keypass confluent -noprompt

echo "SSL credentials generation complete. Files are in the kafka-creds directory."
ls -l

echo "Contents of kafka.client.keystore.jks:"
keytool -list -v -keystore kafka.client.keystore.jks -storepass confluent