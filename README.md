# sales-csv-to-json-api

To start running the project, clone this project, `cd` to the directory and do the followings:

```
docker-compose up -d
docker exec -it mongodb bash
 - mongo
 - use salesDB
 - db.createCollection("sales_records")
```

Now the program is UP and listening to localhost port 3000!

To quickly test the program, run the following:
```
curl -F 'file=/SOME_PATH/sales_json_to_csv_api/testing_files/small.csv' http://localhost:3000/sales/record
```
where SOME_PATH is the path to the cloned repo
