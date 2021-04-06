# sales-csv-to-json-api

To start running the project, clone this project, `cd` to the directory and run the following:

```
docker-compose up -d
```

Now the program is UP and listening to localhost port 3000!

To quickly try the program, run the following:
```
curl -F 'file=/SOME_PATH/sales-csv-to-json-api/testing_files/small.csv' http://localhost:3000/sales/record
```
where SOME_PATH is the path to the cloned repo

To run unit tests, run:
```
make
```
