import { AddAlbums } from "../controllers/album_controller";
import { GetAllAlbum } from "../controllers/album_controller";
import { GetAlbumById } from "../controllers/album_controller";
import { UpdateAlbum } from "../controllers/album_controller";
import { DeleteAlbum } from "../controllers/album_controller";
import { AuthenticateUser } from "../middleware/auth";
import { Router } from "express";

const routerAlbum =  Router();

routerAlbum.post('/add', AuthenticateUser, AddAlbums)
routerAlbum.get('/all', AuthenticateUser, GetAllAlbum)
routerAlbum.get('/get/:id', AuthenticateUser, GetAlbumById)
routerAlbum.put('/update/:id', AuthenticateUser, UpdateAlbum)
routerAlbum.delete('/delete/:id', AuthenticateUser, DeleteAlbum)
export default routerAlbum