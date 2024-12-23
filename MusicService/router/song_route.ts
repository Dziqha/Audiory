import { AddSongs } from "../controllers/song_controller";
import { GetAllSongs } from "../controllers/song_controller";
import { GetSongById } from "../controllers/song_controller";
import { UpdateSong } from "../controllers/song_controller";
import { DeleteSong } from "../controllers/song_controller";
import { AuthenticateAdmin } from "../middleware/auth";
import { Router } from "express";
const routerSong =  Router();

routerSong.post('/add', AuthenticateAdmin, AddSongs)
routerSong.get('/all', GetAllSongs)
routerSong.get('/get/:id', GetSongById)
routerSong.put('/update/:id', AuthenticateAdmin, UpdateSong)
routerSong.delete('/delete/:id', AuthenticateAdmin, DeleteSong)
export default routerSong