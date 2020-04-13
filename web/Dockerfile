FROM nodejs:8.9.4

ADD . /data/webapps
RUN cd /data/webapps && yarn install && yarn build
EXPOSE 8000
RUN chmod 777 run.sh
WORKDIR /data/webapps

CMD ["/data/webapps/run.sh"]
