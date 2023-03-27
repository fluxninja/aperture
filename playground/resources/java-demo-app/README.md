# JavaDemoApp - README

This document explains how to run the JavaDemoApp project from the command line
or an IDEE. The project is using Gradle.

## Requirements

Java Development Kit (JDK) 11 or later Gradle 7.0 or later

## Build and Run

To build the project, run the following command from the root directory of the
project:

`./gradlew build`

This command will download the required dependencies, compile the source code,
and create a runnable jar file in the build/libs/ directory.

To run the project, execute the following command from the root directory of the
project. By default, the application will run on port 8080.

`gradle bootRun`

Once the application is running, you can access it in a web browser by
navigating to http://localhost:<port>, where <port> is the port number on which
the application is running.

## Running multiple instances

To run multiple instances of the JavaDemoApp on different ports, you can use the
following command:

`gradle bootRun --args='--server.port=8080' & gradle bootRun --args='--server.port=8081'`

## Shutdown

To safely bring down both instances, you can use the following command:

`gradle --stop`
