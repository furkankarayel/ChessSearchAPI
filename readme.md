# ChessSearchAPI
<a href="https://postimg.cc/2Lh1XLj5">
  <img src="https://i.postimg.cc/XqD9q9T9/chessearchapi.webp" alt="ChessSearchAPI" width="400"/>
</a>

## Description
This project is a Golang-powered API designed to efficiently query large chess game databases stored in PGN (Portable Game Notation) format. By integrating with powerful chess tools like Scoutfish and Pgn-extract, the API provides a flexible and performant interface for querying games based on various criteria, such as player names, board positions (FEN), or game year. The API is optimized for scalability and precision, making it suitable for chess data analysis, historical research, and other chess-related applications. Development is ongoing, with continuous improvements to extend functionality and performance.

## Dockerfile Explanation

For an efficient docker deployment it's necessary to split up the stages into two parts, build and run. The purpose of it is to leave behind everything you don't need when deploying the ChessSearchAPI. 

However there are four necessary steps in the build stage (Dockerfile has good comments). In summary compiling the ChessParserAPI itself, aswell as Pgn-parser and Scoutfish.

For being able to build both executables, we assign them the right prebuilt docker images that provide you the needed version of your programming language or compiling tools to deploy your apps e.g.: `golang:1.22` and `gcc:14`.

So the build stages are named `go-builder` and `cpp-builder`, which has the purpose of adressability (to copy the binaries/executables) of them later in the run stage.

The run stage is set for running the ChessSearchAPI in Go, yes only for this purpose (as Scoutfish gets executed only on desire). Before we can start the applications in this stage, they need to get copied to the run stage first. Finally we EXPOSE the ports 8080 to the outside world.

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
