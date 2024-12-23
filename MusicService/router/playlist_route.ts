import { AddPlaylist } from "../controllers/playlist_controller";
import { GetAllPlaylists } from "../controllers/playlist_controller";
import { GetPlaylistsById } from "../controllers/playlist_controller";
import { UpdatePlaylist } from "../controllers/playlist_controller";
import { DeletePlaylist } from "../controllers/playlist_controller";
import { AuthenticateUser } from "../middleware/auth";
import { Router } from "express";
const routerPlaylist = Router();


routerPlaylist.post('/add', AuthenticateUser, AddPlaylist)
routerPlaylist.get('/all', AuthenticateUser, GetAllPlaylists)
routerPlaylist.get('/get/:id', AuthenticateUser, GetPlaylistsById)
routerPlaylist.put('/update/:id', AuthenticateUser, UpdatePlaylist)
routerPlaylist.delete('/delete/:id', AuthenticateUser, DeletePlaylist)

export default routerPlaylist;