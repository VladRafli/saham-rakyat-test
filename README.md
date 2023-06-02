# Saham Raykat Backend Test

This repository is used for submission of Backend Test by Saham Rakyat for me. This should my first project using Go Language for creating API.

## Project Requirements

- [X] Please use Golang as a programing language
- [X] Use echo framework
- [X] 5 endpoint (create, list, detail, update, delete) for each table using clean architecture, 
- [X] Use concurrent request to handle multiple request at same time, when insert, get data
- [X] Use Persistent database postgresql/mysql, use redis as cache
- [X] Please show Logs in file, pre request, and post request
- [X] Use environment file (.env) to handle all credentials setting (port, db cred, redis cred )
- [X] Api Documentation (please using postman documentation)
- [X] Put your code and documentation in your Github repository, and please share the url
Repository

## Project Bonus Requirements

- [X] Use Gorm, to use the database
- [X] Use pagination to get list
- [] Build wiith docker and container
- [] Create CI/CD File for deploy from docker
- [] Create Unit test for every service in every layer
- [X] Why clean architecture is good for your application ? if yes, please explain
- [X] How to scale up your application and when it needs to be

## API Documentation

Available to be imported to Postman on root folder. 

## Architecture

This project implements feature-based architecure for more simplified project structure and focused per feature development.

## Things that dont work properly or work not as intended

- Logger only able to log request only.
- Logger need /logs folder present in project in order to work correctly.
- Redis caching might work or not work as on testing with postman, there no significant response time improvement.
- Validation response is hard to manage, so its send as raw string.

## Why using clean architecture

Clean architecture helps on fasten development and structure project more neatly. More well structured project, more easier and faster to work on spesific parts of project.

Clean architecrue also helps to seperate concerns of each part of the project to achieve maximum decoupling and easier to replace parts of the project if there is any breaking changes.

## How to scale up application

When a project need to scale up, surely one that needs attention is to prepare the code that can sustain for future improvement. Thats can be achieved by decoupling the project to its maximum capabilities and prepare the code that can be extended when it needed.

By decoupling the app, we can break down each part of the project, then we add new features or library to project with out worry of breaking old codes.

By prepating the code that can be extended when it needed, we can add new features without touching old codes, but still follows the same flow as the old code. Example of this can be seen when an API needs to be upgraded with breaking changes, we can implement API Versioning so to keep old API working whilst working on new API. 

## License

The Unlicense