# ChessSearchAPI
## Description
This is an API which is responsible for querying pgn databases of chess games efficiently. It provides an interface to Scoutfish and Pgn-extract, which allow querying based on different criteria.

## Dockerfile Explanation

For an efficient docker deployment it's necessary to split up the stages into two parts, build and run. The purpose of it is to leave behind everything you don't need when deploying the parser service. 

However there are two necessary steps in the build stage. Compiling the go parser service itself and currently the c++ parser called Scoutfish (pgn-parser will be added soon).

For being able to build both executables, we assign them the right prebuilt docker images that provide you the needed version of your programming language or compiling tools to deploy your apps e.g.: `golang:1.18` and `gcc:14`.

So the build stages are named `go-builder` and `cpp-builder`, which has the purpose of adressability (to copy the binaries/executables) of them later in the run stage.

The run stage is set for running the parser service in Go, yes only for this purpose (as Scoutfish gets executed only on desire). Before we can start the applications in this stage, they need to get copied to the run stage first. Finally we EXPOSE the ports 8080 to the outside world.

The comments in the Dockerfile could also help to understand the context.

## Docker Run
For building the docker container it's necessary to have docker installed for CLI, you should be able to run ```docker --version```

running the command below builds our image (that runs through both stages and deploys our parser-service) and gives it a tag name, else we would have to find out the random id that's given to our built container 

`docker build --no-cache -t parser-service .`

to run it we run the following command which maps the host port to container port

`docker run -p 8080:8080 -it parser-service`

After this go into your browser type localhost:8080 and you should be able to see a message of our service.

## Dockerfile Notes

`WORKDIR` = defines where the files are located within that specific STAGE

`COPY --from=<STAGE_NAME> <LOCATION_IN_PREVIOUS_STAGE> <LOCATION_IN_CURRENT_STAGE> `= this command is used to copy our binaries from the build stages into the run stage, and also the files we need for execution like the pgn folder for the GigaBase scout file

`RUN` = a command that makes you execute whatever you want 

`EXPOSE <PORT_NUMBER>` = let's you predefine which port is going to be exposed 
