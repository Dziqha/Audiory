import { AddArtist } from "../controllers/artist_controller";
import { GetALlArtist } from "../controllers/artist_controller";
import { GetArtistById } from "../controllers/artist_controller";
import { UpdateArtist } from "../controllers/artist_controller";
import { DeleteArtist } from "../controllers/artist_controller";
import { AuthenticateAdmin } from "../middleware/auth";
import { Router } from "express";

const routerArtist = Router();


routerArtist.post("/add", AuthenticateAdmin, AddArtist)
routerArtist.get("/all", GetALlArtist)
routerArtist.post("/get/:id",GetArtistById )
routerArtist.put("/update/:id", AuthenticateAdmin, UpdateArtist)
routerArtist.delete("/delete/:id", AuthenticateAdmin, DeleteArtist)
export default routerArtist