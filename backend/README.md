# Backend service for MeetOver
Talks to REST client and firebase
## REST
Gorilla mux for rest router

## Firebase
uses schema ...
stored in a local JSON object

## URLs available:
    - location/<coordinate string>
    - match/<other user's id>
    - people/<id of person to get>
    - login/<credential>

## Deployment on Heroku
Some steps here...
``git subtree push --prefix backend heroku master``

## Installation instructions
Place current directory in GOPATH

``go build ``

``backend``

Access the REST service using the URLs on the host machine's IP address and port 8080.