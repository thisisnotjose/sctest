# Follower-Maze Solution

This is a solution of the follower-maze challenge provided by `SoundCloud` implemented in `Go`. If you need information on the specifications of the challenge please go to the [instructions](./docs/instructions.md)

## Getting Started

In order to execute the program you can do: 

```
$ make start
```

Or if you want to use docker

```
$ docker-compose up
```

## Testing 

To test the application you can do

```
$ make test
```

Or if you want to use docker
```
docker-compose run --rm sctest make test
```

## Part 1 (see [PR #1](https://github.com/thisisnotjose/sctest/pull/1))

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

- **Unit tests**: Tests are still not in the code, but the different components should make it easy to mock any dependencies and test the smallest ammounts of code possible at a time. Also it would be great to add tests to the overall functionality and not just the modules, meaning having tests that boot the application up and then mock the clients, connections and requests done to the service.
- **Logging**: Although there is logging in the application, the use of LEVEL(INFO, DEBUG, WARN, ERR) logging would be great.
- **Environments**: The use of environment to define basic behaviour and/or configurations of the app, as is the use of environment variables to load things like the ports.
- **Persistence**: It would be valuable to move things like the follow registry to an actual persisted layer, like Redis or Mongo or Postgres. Depends a lot of the production requirements or expected usage.
- **Use more than one process**: For this we could set `GOMAXPROCS` to the number of go routines we have and it would allow every go routine to take advantage of using its own processing unit. The only thing that we would need to do is make the `context` thread safe.
- **Decouple the processors**: If we wanted to distribute this it could change the channel for a rabbitMQ queue and decouple the event processor into its own service.
- **Abstract the logic for sending a client message**: There should a single piece of code that is used across the application to send a message to a client instead of making `fmt` calls in each event processor call.
- **Dockerize the application**: This can be done with minimal amount of work and it could simplify the way its deployed and tested.

## Part 2 (see PR [PR #2](https://github.com/thisisnotjose/sctest/pull/2))

The second part is mostly dedicated to the Dead Letter Queue implementation. Several modifications to the code were done in order to track and preserve the messages coming to the events queue that couldn't be processed for one or another reason. 

### Approach

The approach taken is inline with the event queue processing, a channel for the dead letter queue was created and a processor initiated when booting up the application, this processor waits for message strings coming into the queue (they needed to be strings and not event types because in several cases the string is malformed and cannot be unmarshalled) and then prints them out to the console.

The idea of using a channel and moving this out of the logic of the main routine is to avoid spending time when processing events on processing dead letters, we should as much as possible be delegating that to another process and then using the data gathered in that other process to understand how is the usage of the events layer.


### DLQ Treatment

There are several actions we can take with this queue:

- Setup a dashboard that shows the dead letter count over time, the types of dead letters and their count (malformed, incomplete, no user connected) 
- Depending on the context it might be valuable to setup some alarms like "if we receive more than 5 dead letters in an hour trigger alarm" or "if we receive a malformed dead letter raise an alarm"
- Do a weekly report where it shows the dead letter count per application, that way its visible to other developers because they might have a bug and not be aware of it.
- If there are too many dead letters for unconnected users we should evaluate the code, are we storing the connections correctly? Are we taking too much or taking to little time to send the notification back? Does it even make sense to store this type of exceptions on dead letter queues? 

### Shortcuts and Design tradeoffs

By adding another channel and another processor we are increasing the memory and decreasing the overall performance of the application(same resources shared by a bigger amount of routines), BUT, the big trade-off here is the ability to decouple this functionality in a clear way, which would start by moving the channels to some type of message queue and moving the code of the processors into their own services.

The solution also creates the `users` package and moves the functionality of sending of the events to other users to that package, however, because of time constraints I didn't move the follow/unfollow behavior to that package even though they belong there. 

Another shortcut was that for `broadcast` I'm not adding the messages to the deadletter queue. I could easily add that, but (and this would be something I would ask the product manager) it seems like broadcast is oriented to all `connected` users and not to all users, so there aren't any dead letters in that sense.

And last but not least, when failing to write to a TCP connection of a user I'm not sending the error to the dead letter queue, I'm just logging it, ideally I would put them together to also give context of why the event is inside the dead letter queue. 

### Changes 

- Added a `DeadLetter` channel to the context.
- Whenever the message cannot be processed add it to the channel, this includes:
  - When the event type is not recognized
  - When the message cannot be parsed into an event object because of lack of pipes
  - Or cannot be parsed because of letters where there should only be digits
  - When the message could not be delivered to the user
- Added unit tests for the dead letter queue and some other scenarios related to the dead letter code.
- Added the `users` package with methods on how to send events to users. 
