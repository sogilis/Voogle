# ADR-011-live-update

* Creation Date: 18/05/22

## Context

While developing the front-end environment, we decided to reflect in real time the encoding of the video to the user.
We needed a way to get the information from back to front.

## Problematic

The status front-side need to be updated at the same time than the database if possible.
We may have to check the status of multiple videos at a time and maybe for multiple users as well.

## Decision

We choose to use Websockets to adress our problem. Since Voogle is a way for us to learn, we deemed the use of Websocket to be a good exercice. Plus, as Voogle is a demonstrator, we want to reduce the strain on the server side and reflect the changes made in the database the fastest way possible.

## Options

### Option 1: Polling

The simplest way to get information for a video. We send a request to the server who respond with the status of the video.

#### Benefits

* Simple to set and use.

#### Drawbacks

* We need to repeat the request until the video is fully uploaded.
* The interval between request need to be low to reflect change in the database when they happens.
* The repeat for each video and for each user can cause a strain on both front and back end.

### Option 2: Websocket

Opening a channel between the client and the server allow for bidirectional communication. The server can then choose when to send message.

#### Benefits

* Allow to open a single websocket per user.
* Server can send update to client when needed.

#### Drawbacks

* More difficult to set up.
* User authentification should be adressed in a different way.

### Options diff

Polling has the advantage of being easier to set up, but can put too much of a strain on the server when the number of video to follow increase.
Websocket are way harder to set up, but should reduce the amount of message from server drastically.