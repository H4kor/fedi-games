Fedi-Games is a simple framework to build games for the [Fediverse](https://en.wikipedia.org/wiki/Fediverse) using the [ActivityPub](https://activitypub.rocks/) protocol.

Fedi Game instances:
- [games.rerere.org](https://games.rerere.org)


### Configuration

Fedi Games is configured via environment variables.
These configurations should be set otherwise defaults for local development are used.

**Example:**
```
FEDI_GAMES_HOST=games.rerere.org
FEDI_GAMES_PROTOCOL=https
DATABASE_URL=/path/to/sqlite.db
MEDIA_PATH=/path/to/media/
```

### Development

For development an [air](https://github.com/cosmtrek/air) configurations is provided and the project can be started with `air` to run the server with live reload.

Otherwise use:

```
go run rerere.org/fedi-games/cmd
```

The server will run on [http://localhost:4040](http://localhost:4040). 
In local development (when the `FEDI_GAMES_HOST` env var is not set or inlcudes "localhost") no messages will be send to other servers of the fediverse.

### Building a game

Every game has to implement the `games.Game` interface. 
The `NewState()` function must return a pointer to a JSON serializable struct.
The state doesn't have to be initialized.
On the first message for a new game the game will be passed a zero valued instance of the struct.

The `OnMsg` function will be called for each message sent to the game. 
It must return:
- the updated state of the game
- a `games.GameReply` object to be send as Note via ActivityPub
- an error if the message could not be processed. An error must not be returned for user input errors. Sent a message to the player instead

#### Adding images and attachments

To add images or other attachments to a reply the media service should be used.
By calling `internal.StoreMedia` you can store a blob of data, which will be served on the media route.
This can be added as an attachment to the `GameReply` object.

**Example:**
```go
imgUrl, err := internal.StoreMedia(buffer.Bytes(), "png")
reply := games.GameReply{
    Attachments: []games.GameAttachment{
        {
            Url:       imgUrl,
            MediaType: "image/png",
        },
    },
}
```

