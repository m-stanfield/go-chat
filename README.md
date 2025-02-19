# Project go-chat 

A side project I wanted to do to learn a variety of topics better. The overall goal of this would be to create a "chat app" similar to Teams or Discord, where users could send and recieve messages in different "chat rooms" that are user defined.

The technologies that I wanted to explore are: 
* Websockets
* React
* Golang HTTP Server
  - Specifically the new routing introduced in 1.22

The goal of this project would be to have the following features:
1) Users be able to create accounts
2) Users be able to create servers, which are a collection of chat channels
3) Users be able to create chanenls within a server
4) Validate user has access to an indvidiual server
5) Have method to invite new users to a server
6) Notify users when messages are recieved


# How to build

Due to being under development, build process, dependencies, and documentation has not been exaustively tested on multiple computers. This will be completed later once the program is in a more releaseable state.

Current known dependencies:
* Go v1.23
* Node 18

To run, the `make run` command in the base directory should work.
