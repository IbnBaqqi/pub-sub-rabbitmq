# rabbitmq using stomp protocol
# build with: docker build -t rabbitmq-stomp .
# to run container and management ui
# run with: docker run -d --name stomp_rabbitmq -p 61613:61613 -p 15672:15672 rabbitmq-stomp
FROM rabbitmq:4-management
RUN rabbitmq-plugins enable rabbitmq_stomp