Apache Kafka
=================
See [web site](https://kafka.apache.org) for more details about kafka.

## **Architecture Diagram**

![Architecture Diagram](https://miro.medium.com/max/1400/1*0QqLTumYuNrpmNoZr1QoEA.png)

### **Staging Server**
Please check out kafka server [here](http://15.206.67.39:9000/).

- We hosted our staging kafka server in an ec2 instance of region ap-south-1.
- Kafdrop which is GUI for kafks is also self hosted along with the kafka server.
- Kakfa topics can created using the kafdrop GUI.
- We need to update our .env file with `KAFKA_HOSTS=15.206.67.39:9091,`. Kafka hosts are the kafka brokers here and we can have multiple brokers.
So we are considering the kafka hosts separated by comma.
