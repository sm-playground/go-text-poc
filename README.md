# go-text-poc - REST API with Golang, Mux, Postges (gorm), and Redis with connection pool

Demonstrates golang REST API implementation connecting to the Postgres database using gorm (ORM).

The docker cheat sheet is available here -> https://dockerlabs.collabnix.com/docker/cheatsheet/

1. Go and Docker are installed and configured.
2. To create a docker container running an instance of PostgeSql run the command below:  

-- This command will create the the postgres container with no data.  
docker run --rm -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ia_text --name my-postgres -p 54320:5432 postgres

-- The following command will create a postgres docker container preserving the data locally   
docker run --rm -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ia_text --name my-postgres -p 54320:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres


--rm instucts to remove container after it exists.  
POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB are environment variables specifying parameters of the postgres instance.  
my-postgres is the name of the container  
-p 54320:5432 - connects the port 5432 inside Docker as port 54320 on the host machine. 
-v $HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres - all the data is stored and kept in local folder $HOME/docker/volumes/postgres.   

3. Run psql - an interactive terminal to work with the postgres database.  
docker exec -it my-postgres psql -U postgres. 

4. Useful psql commands:  
- \t - list all databases
- \c database-name - switch the database
- \d \d+ - list of all tables in the database
- \d \d+ table-name - describe the table
  
5. Create a docker container for running Redis image.   
docker run --name my-redis -p 6379:6379 -d redis. 
 
6. To work with the redis cli need to know the container id. The command below will list all running docker containers.   
docker ps. 

7. Access the container's shell by running the following command and providing the correct container id  
docker exec -it 103a71bf2cb3 sh

8. Once in the shell run the cli tool for redis.   
redis-cli
