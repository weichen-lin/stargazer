FROM docker.elastic.co/logstash/logstash:8.12.2

# delete default configuration
RUN rm -f /usr/share/logstash/pipeline/logstash.conf

# copy configurations
COPY pipeline/ /usr/share/logstash/pipeline/
COPY config/ /usr/share/logstash/config/