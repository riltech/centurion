# Architecture

## The overarching idea

Centurion uses a domain based architecture. It is utilising multiple popular patterns and principles to create a scalable and easy to understand codebase.

One of the core principles is the [single responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle). All packages are structured based on their responsibility in the system. Every package is individual and can be used from any other package in an upper level of abstraction.

As there are two easily identifiable main parts of the system:
* engine - Responsible for the game mechanics and communication
* dashboard - Visual feedback about the game state

Centurion needs a way to provide communication between these two packages. As they are in the same layer it would be unwise if this communication is direct, therefore we are implementing an event bus for them which they can use to communicate between each other.

This event bus also provides a very easy way to introduce more and more business logic if we want to as any new package in this layer could work with the already defined events of the event bus therefore allowing the system to be almost infinitaly scalable in terms of code extension.

In the following packages you can see very distinctive __"controller - service - repository - domain model"__ pattern:
* challenge
* player
* combat
* engine

Inside this packages this pattern allows us to completely separate abstraction layers further and to understand the nature of a given package. While `engine` is different than the other packages (mostly because of not having a repository and model), but we can see that the same abstraction strategy works just fine with this package as well. 

## Package roles

#### package main

```sh
core/ 
main.go
```
The root is where bootstrapping and the configuration setup should happen. I also consider it to be a good pattern to try orchestrating a graceful shutdown from this package.

#### package core

```sh
core/
## Packages
 - bus/
 - challenge/
 - combat/
 - config/
 - dashboard/
 - engine/
 - player/
 ## Files
 - dashboard.go
 - engine.go
 - exit.go
```

In this layer we are still having bootstrapping logic. In our case we are running an `engine` for the game that handles http and websocket traffic. As it should be easy to follow the current state of the game we have a `dashboard` which is the visualisation module. Our `exit.go` is responsible for providing a handy interface to orchestrate graceful exit.

#### package bus

```sh
core/bus/
 ## Files
 - bus.go
 - events.go
```

Event bus implementation that provides an easy to use interface for communication. Under the hood this packages utilises unbuffered channels. The `events.go` file holds all the events that can happen over the bus. This is important because of encapsulation, as in this way the bus package has no other dependencies in the system, which means that it can be used from anywhere as it is a standalone package.

#### package challenge

```sh
core/challenge/
 ## Files
 - challenge.go
 - default.go
 - repository.go
 - service.go
```

Challenge is a typical domain module. `challenge.go` holds all the model information that describe a challenge. We also store the enum values for challenge types here. `default.go` describes the default challenge modules that are part of the system. `repository.go` is used for implementing storage and last but not least the `service.go` exposes storage and business functionalites.
