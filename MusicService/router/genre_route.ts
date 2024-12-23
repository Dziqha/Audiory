import { AddGenre } from "../controllers/genre_controller";
import { GetALlGenre } from "../controllers/genre_controller";
import { GetGenreById } from "../controllers/genre_controller";
import { UpdateGenre } from "../controllers/genre_controller";
import { DeleteGenre } from "../controllers/genre_controller";
import { AuthenticateAdmin } from "../middleware/auth";
import { Router } from "express";

const routerGenre = Router();

routerGenre.post("/add", AuthenticateAdmin, AddGenre)
routerGenre.get("/all", GetALlGenre)
routerGenre.post("/get/:id",GetGenreById )
routerGenre.put("/update/:id", AuthenticateAdmin, UpdateGenre)
routerGenre.delete("/delete/:id", AuthenticateAdmin, DeleteGenre)

export default routerGenre