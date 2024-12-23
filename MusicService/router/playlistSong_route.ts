import { AddPlaylistSong } from "../controllers/playlistSong_controller";
import { GetAllPlaylistsSongs } from "../controllers/playlistSong_controller";
import { GetPlaylistSongById } from "../controllers/playlistSong_controller";
import { UpdatePlaylistSong } from "../controllers/playlistSong_controller";
import { DeletePlaylistSong } from "../controllers/playlistSong_controller";
import { AuthenticateUser } from "../middleware/auth";
import { Router } from "express";

const routerPlaylistSong = Router();

routerPlaylistSong.post('/add', AuthenticateUser, AddPlaylistSong)
routerPlaylistSong.get('/all', AuthenticateUser, GetAllPlaylistsSongs)
routerPlaylistSong.get('/get/:id', AuthenticateUser, GetPlaylistSongById)
routerPlaylistSong.put('/update/:id', AuthenticateUser, UpdatePlaylistSong)
routerPlaylistSong.delete('/delete/:id', AuthenticateUser, DeletePlaylistSong)

export default routerPlaylistSong;