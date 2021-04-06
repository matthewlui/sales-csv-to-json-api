# sales-csv-to-json-api
<b>Prerequisite:</b></br>
The host machine should have Docker and GNU Make installed.</br></br>

<b>1. To start running the project, clone this repo, `cd` to the directory and run the following:</b>

```
make start
```

Now the program is UP and listening to localhost port 3000!</br></br>


<b>2. To upload a csv, run the following:</b>
```
curl -F 'file=/SOME_PATH/sales-csv-to-json-api/testing_files/small.csv' http://localhost:3000/sales/record
```
where SOME_PATH is the path of the cloned repo.</br></br>
 
<b>3. To query the result, open internet browser and go to:</b>
```
  http://localhost:3000/sales/report
```
</br></br>
<b>4. To run unit tests, run:</b>
```
make test
```
