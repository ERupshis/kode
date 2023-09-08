# test task for KODE
description:
It is necessary to design and implement a service in Golang, providing REST API interface with methods:
- add note;
- output the list of notes.

Data should be stored in PostgreSQL.
When saving notes it is necessary to validate spelling errors using the Yandex.Speller service (add integration with the service).

It is also necessary to implement authentication and authorization. Users should have access only to their notes. The possibility of registration is not mandatory, it is acceptable to have a predefined set of users (the mechanism of storing accounts is any, up to hardcode in the application).


# Comments for checking:
hardcoded users:
- asd : asd
- qwe : qwe

Docker:
- split on two containers: service and db
- steps for launch:
  - docker build -t kode:multistage -f Dockerfile.multistage .
  - docker compose config
  - docker compose up --build

Postman collection is located in root repository:
- https://github.com/ERupshis/kode/blob/864b7ca89bf3e490a6476e2c464d21387ed20c35/KODE%20test%20task.postman_collection.json

Curl examples are located in root repository:
- https://github.com/ERupshis/kode/blob/69f53c96b63ffff38cad09f0e1e73f5b61e41201/curl_examples.txt
