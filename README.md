# Follower-Maze Solution

This is a solution of the follower-maze challenge provided by `SoundCloud` implemented in `Go`. If you need information on the specifications of the challenge please go to the [instructions](docs/instructions.md)

## Getting Started

In order to execute the program you can do: 

```
$ make start
```

Or if you want to use dockerq

```
$ docker-compose up
```

## Part 1

During this part the code was refactored in the following manner: 

- The main was moved to cmd/main
- a `MakeFile` was created
- a `docker-compose` configuration was added.
- Added the concept of `servers`, as the abstraction of the TCP listener for clients and servers
- Added the concept of `handlers` as the method that gets executed in each message coming to the server 
- Astracted the specific logic of each type of event to the `events` package
- Created the concept of `processor` as a worker that fetches from a queue and executes the event
- Added a `channel` for filling the event queue

### Code Overview

The following will go over technologies used, project structure, main components and possible improvements to give a comprensive description of the approach used and why.

I also suggest reviewing the code, since the functionality and decisions are also explained inline.

#### Technologies

For the code language I chose `Go`, mainly because I think its standard libraries and behaviours, in concurrency(go routines) and thread communication(channels) would give an advantage when working with this specific problem.

#### Project Structure

Ideally the structure was changed in order to abstract the behaviour of the technology, meaning the TCP channels breaking the messages with pipes and control characters vs the actual logic of the application, the events and their actions:

  - **cmd/main.go**: Main file, this boots the application by creating subroutines and maintaining the context object.
  - **internal**: All logic files directly related to the functionality of this app.
    - **servers**: This are the abstraction classes for the net libraries, they listen to the socket, read the information and abstract as much functionality as possible in order to only send a string message of the body of the event to the next component.
    - **handlers**: This would be the components dedicated to processing the event, in `MVC` this would be the closest to a `controller`. 
    - **events**: This are helper methods that describe the action of an event, so there is one file for `follow`, `unfollow`, `private message`, `status update` and `broadcast`. This way there is an obvious place to look at the specific behaviour of an event. Most of the code is general, but this is pretty specific to the logic of the application so is heavely abstracted to make it easily testable.
    - **types**: Are the different types of classes used across the application, from `Event` to `Server` or `Processor`. This not only allows to have clear interfaces across the code but to potentially use this interfaces to mock the other components in the application when using unit tests.
  - **docs**: The documentation provided by soundcloud.
  - **original**: The original solution provided by sloundcloud
  - **Makefile**: The makefile with instructions for running the application.
  - **docker-compose**: A development docker compose file to create an instance with go and test the application.

#### Biggest Changes

- **Use Go routines**: Divide the code into three main functions, listen client connections, listen to events and processing the events. For this we separate the logic from servers and processors, each of them its own subroutine.
- **Context**: As a way to maintain the same context for all three routines we create a `context` object and pass it as reference to each routine. 
- **Channels to communicate the event listener and the processor**: The processor is in charge of adding events to the queue and processing the events that come, but if idle the server should not be checking for new messages constantly, so we created a channel between the handler and the event processor and pass along the information through there.
- **Handlers and Servers**: Theres is quite a bit of logic that can be abstracted as middleware, this solution is not quite there yet, but at least does a first level abstraction of the code that is related to the messaging protocol vs the business logic.


#### Possible Improvements

- **Unit tests**: Tests are still not in the code, but the different components should make it easy to mock any dependencies and test the smallest ammounts of code possible at a time.
- **Logging**: Although there is logging in the application, the use of LEVEL(INFO, DEBUG, WARN, ERR) logging would be great.
- **Environments**: The use of environment to define basic behaviour and/or configurations of the app, as is the use of environment variables to load things like the ports.
- **Persistence**: It would be valuable to move things like the follow registry to an actual persisted layer, like Redis or Mongo or Postgres. Depends a lot of the production requirements or expected usage.
- **Decouple the event processor**: If we wanted to distribute this it could change the channel for a rabbitMQ queue and decouple the event processor into its own service.
- **Abstract the logic for sending a client message**: There should a single piece of code that is used across the application to send a message to a client instead of making `fmt` calls in each event processor call.

