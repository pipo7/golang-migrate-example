# Db Migration process
- Apply SQL DDL and DML 
- Ensure data is valid and ETL process between source and target
- Ensure the Rollback exists  
    - Either as Rollback SQL scripts OR
    - take backup of DB and restore it in case of rollback

# Commonly used libraries for Db Migration
Flyway , Liquibase , golang-migrate , Goose
Only in  Liquibase we have (For complex Dbs, multiversion support and order of execuxn of cmds can be changed )
Rest all 3 are linear DB versioning support, cant change order .

# Run postgres container - direct Cmd
```sudo docker run --rm --name hostman-pgdocker -e POSTGRES_PASSWORD=hostman -e POSTGRES_USER=hostman -e POSTGRES_DB=hostman -d -p 5531:5531 -v $HOME/docker/volumes/postgres:/var/lib/ postgresql/datapostgres```

--rm tells the system to delete the container and its file system after the container is stopped. It helps saving server space.
--name is the container name, which must be unique within one server, regardless of status. here it is ```hostman-pgdocker```
-e points to the environment variables: name and password of the superuser, default database name.
-d launches the container in background (offline) mode. After the launch, control is returned to the user.
-p binds the Postgres port to the server port.
-v creates a mount point.

# OR just Use docker-compose.yaml file and run
``docker compose up -d``
To disconnect and remove all volumes try 
``docker compose down -v``

# Connecting to PgSQL using psql OR docker exec

connect to the isolated environment using the psql utility. It is included in the postgresql-client. Run:
``sudo apt install postgresql-client``
Connect to the database and execute a test query:
``psql -h 127.0.0.1 -U hostman -d hostman``
The output should be:
Password for user hostman:
psql (12.6 (Ubuntu 12.6-0ubuntu0.20.04.1), server 13.2 (Debian 13.2-1.pgdg100+1))
Type "help" for help.
hostman=#

OR better use docker exec 
```sudo docker exec -it hostman-pgdocker psql -U hostman```

# Example to create a Table
``docker ps`` 
CONTAINER ID   IMAGE           COMMAND                  CREATED          STATUS                    PORTS                                                 NAMES
0a8a0a2fb343   postgres:13.3   "docker-entrypoint.s…"   20 minutes ago   Up 20 minutes (healthy)   5432/tcp, 0.0.0.0:5531->5531/tcp, :::5531->5531/tcp   postgres_migration_go-postgres-1

``docker exec -it postgres_migration_go-postgres-1 psql -U testuser -d testdb``
psql (13.3 (Debian 13.3-1.pgdg100+1))
Type "help" for help.
testdb=# create table cities (name varchar(80));
CREATE TABLE
testdb=# insert into cities values ('Seattle');
INSERT 0 1
testdb=# select * from cities;
  name   
---------
 Seattle
(1 row)


# Migrations
First we apply DOWN in desc order of numbers i.e. 3000 then 2000 and then 1000
Then we apply UP in asc order - 1000 -> 2000 -> 3000

File *.sql format is version_<description>.down.sql  OR version_<description>.up.sql
During UPGRADE it goes from lowest to highest version 
During DOWNGRADE from highest to lowest version

m.Down followed by m.UP
OR you can direclty go to a certain step also m.Step(2)

_ = m.Steps(2) // NOTE This executes only 2 steps that is versions 1000 and 2000 thus After this
	// statement you see the session_review column is NOT added

_ = m.Force(1000) // NOTE This forces DB schema to be at version 1000
Suppose you get an error in version 2000 then you can go back to version 1000 and correct the next version SQLs . Thus FORCE(version#) forces to rollback till that verison

# NOTE
sudo docker stop hostman-pgdocker
This is where we discover a security hole. The files and processes that the container creates are owned by the internal postgres user. Due to the lack of a namespace within the container, the UID and GID can be arbitrary.

The problem is that there may be other privileged users and groups on the server itself with the same UID and GID. This potentially leads to a situation where the user can access host directories and kill any processes. You can avoid this situation by specifying the USERMAP_UID and USERMAP_GID variables at startup.
```sudo docker run --rm --name hostman-pgdocker -e POSTGRES_PASSWORD=hostman -e POSTGRES_USER=hostman -e POSTGRES_DB=hostman -e USERMAP_UID=999 -e USERMAP_GID=999 -d -p 5432:5432 -v $HOME/docker /volumes/postgres:/var/lib/postgresql/data postgres```


# Run the go main.go to conenct to Db
You may have to change the initial password 
``ALTER USER testuser WITH PASSWORD 'newpassword';``
Change the new password in docker-compose.yaml and in main.go connection string
