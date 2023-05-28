FROM registry.access.redhat.com/ubi8/ubi-minimal
LABEL description="Custom image for devops test"
RUN mkdir /opt/app
COPY ./devops /opt/app
RUN chmod +x /opt/app/devops && microdnf install shadow-utils && adduser \
       --no-create-home \
       --system \
       --shell /usr/sbin/nologin \
       devops
USER devops
EXPOSE 8080
ENTRYPOINT ["/opt/app/devops"]
