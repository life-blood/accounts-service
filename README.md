## Microservice for donations exposing REST-API with MySQL database

Supports CRUD operations for Donors and Acceptors entities

See project docs: https://github.com/life-blood/documentation
## Run locally
 ### Prerequisites
- Docker Desktop installed on your machine

### Installation
```$ git clone https://github.com/life-blood/accounts-service/```

Open a new window in your terminal to start the docker instance hosting the MYSql image and run:

```$ docker-compose up ```

Then start the service:

``` $ go run main.go ```
## LifeBlood Project Architecture
![alt text](https://i.ibb.co/M7C45Wv/Architecture.png)
