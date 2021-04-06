# sales-csv-to-json-api
<b>Prerequisite:</b></br>
The host machine should have Docker installed.</br></br>

To start running the project, clone this repo, `cd` to the directory and run the following:

```
docker-compose up -d
```

Now the program is UP and listening to localhost port 3000!</br></br>


To upload a csv, run the following:
```
curl -F 'file=/SOME_PATH/sales-csv-to-json-api/testing_files/small.csv' http://localhost:3000/sales/record
```
where SOME_PATH is the path of the cloned repo.</br></br>
 
To query the result, open internet browser and go to:
```
  http://localhost:3000/sales/report
```
</br></br>
To run unit tests, run:
```
make
```
