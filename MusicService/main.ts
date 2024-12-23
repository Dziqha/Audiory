import routerArtist from "./router/artist_route";
import routerGenre from "./router/genre_route";
import routerAlbum from "./router/album_route";
import routerSong from "./router/song_route";
import routerPlaylist from "./router/playlist_route";
import routerPlaylistSong from "./router/playlistSong_route";
import routerCollaborator from "./router/collaborator_route";
import routerSongCollaborator from "./router/songCollaborator_route";
import { getSongById } from "./controllers/song_controller";
import { Initialize } from "./configs/redis";
import express from "express";
import cros from "cors";
const app = express();
import * as grpc from '@grpc/grpc-js'
import { loadSync } from "@grpc/proto-loader";
import { getArtistById } from "./controllers/artist_controller";
import { getAlbumById } from "./controllers/album_controller";
const PROTO_PATH = "./song.proto";
const packageDefinition = loadSync(PROTO_PATH, {
  });
  
const songProto = grpc.loadPackageDefinition(packageDefinition).song as any;
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(cros());
app.use("/artist", routerArtist);
app.use("/genre", routerGenre);
app.use("/album", routerAlbum);
app.use("/song", routerSong);
app.use("/playlist", routerPlaylist);
app.use("/playlistSong", routerPlaylistSong);
app.use("/collaborator", routerCollaborator);
app.use("/songCollaborator", routerSongCollaborator);


// const port = process.env.NODE_PRODUCTION || "production";

const port = 3000;
app.listen(port, () => {
    console.log(`🚀 Server ready at Port ${port}`);
});

Initialize();
const grpcServer = new grpc.Server();
grpcServer.addService(
    songProto.SongService.service, {
    GetSongById: getSongById
});

grpcServer.addService(
    songProto.ArtistService.service, {
    GetArtistById: getArtistById
});

grpcServer.addService(
    songProto.AlbumService.service, {
        GetAlbumById: getAlbumById
    }
)

const PORT = 50051;
grpcServer.bindAsync(`0.0.0.0:${PORT}`, grpc.ServerCredentials.createInsecure(), () => {
    console.log(`🚀 GRPC server running at Port ${PORT}`);
});