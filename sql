sql:

CREATE TABLE track_list (
    id        SERIAL NOT NULL PRIMARY KEY,
    name      VARCHAR(60) NOT NULL,
    artist      VARCHAR(60) NOT NULL,
    album      VARCHAR(60) NOT NULL
    );
CREATE TABLE like_list (
    id        SERIAL NOT NULL PRIMARY KEY,
    username  VARCHAR(60) NOT NULL,
    track_id  INTEGER NOT NULL REFERENCES track_list(id) ON UPDATE CASCADE
    );

query:
protobuf.protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative like/like.proto