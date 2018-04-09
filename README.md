# MeetOver
## Description
Mobile application for iOS and Android used to connect busy professionals on the fly

## Backend services
### REST
- Gorilla mux for rest router
- Each route is implemented under `/router/handlers`
- Default port: `8080`

### Firebase
- Realtime Database for user data (via firego)
- Authentication for users
- Storage for large files

### Endpoints:
- Unauthenticated
    - `GET: /`
    - `POST: /login/<linkedin code>`
    - `POST: /test/<type of test>`
- Authenticated
    - `POST: /refreshtoken`
    - `POST: /meetover/<other user id>`
    - `POST: /meetover/decision/<other user id>`
    - `POST: /userprofiles`
    - `POST: /sendpush`
    - `POST: /match/<other user id>`

### Installation instructions
1. Place current directory in `GOPATH`
1. Navigate to `/backend`
1. Install dependencies: `dep ensure`
1. Build: `go build`
1. Run: `./backend`

### Deployment on Heroku
Command: `git subtree push --prefix backend heroku master`

## Client App
### Stack
#### React Native
- Core native libraries

#### Expo
- Seamless development cycle
- Push notification support
- Geo location support

#### Redux
- State management
- Allows for state access from any component

#### React Navigation Routes
- `LoginScreen`
- `CreateProfileScreen`
- `AppNavigator`
    - `MainTabNavigator`
        - `ListScreen`
        - `MapScreen`
        - `ChatsScreen`
    - `SettingsScreen`
    - `ProfileScreen`
    - `RequestScreen`
    - `ConfirmScreen`
    - `ChatScreen`

#### Firebase
- Realtime Database
    - User can make authenticated requests to data from client
    - Allows direct integration for chat feature
- Authentication
    - Uses JWT to securely authenticate users

#### NativeBase
- Main source of cross-platform UI components

#### Gifted Chat
- UI chat component for user to user communication

#### Yarn
- Package manager

### Installation instructions
1. `yarn install` or `npm install`
1. `yarn start` or `npm start`
1. Install the Expo app
1. Connect to the development server
