FROM alpine:3.18

ARG TERRAFORM_VERSION=1.9.3
ARG ANSIBLE_VERSION=10.2.0

RUN apk add --no-cache \
    bash \
    curl \
    python3 \
    py3-pip \
    aws-cli \
    && pip3 install --no-cache-dir --upgrade pip

RUN curl -LO "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

RUN curl -L -o /usr/local/bin/opa https://openpolicyagent.org/downloads/latest/opa_linux_amd64_static \
    && chmod 755 /usr/local/bin/opa

RUN pip3 install --no-cache-dir ansible==${ANSIBLE_VERSION}

COPY kado /usr/local/bin/kado

ENV PATH="/opt/venv/bin:$PATH"

WORKDIR /workspace

#ENTRYPOINT ["tail", "-f", "/dev/null"]
ENTRYPOINT ["/usr/local/bin/kado"]