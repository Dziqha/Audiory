import { AddSongCollaborator } from "../controllers/songCollaborator_controller";
import { FindAllSongCollaborator } from "../controllers/songCollaborator_controller";
import { FindSongCollaboratorById } from "../controllers/songCollaborator_controller";
import { UpdateSongCollaborator } from "../controllers/songCollaborator_controller";
import { DeleteSongCollaborator } from "../controllers/songCollaborator_controller";
import { AuthenticateAdmin } from "../middleware/auth";
import { Router } from "express";
const routerSongCollaborator =  Router();

routerSongCollaborator.post('/add', AuthenticateAdmin, AddSongCollaborator)
routerSongCollaborator.get('/all', FindAllSongCollaborator)
routerSongCollaborator.get('/get/:id', FindSongCollaboratorById)
routerSongCollaborator.put('/update/:id', AuthenticateAdmin, UpdateSongCollaborator)
routerSongCollaborator.delete('/delete/:id', AuthenticateAdmin, DeleteSongCollaborator)
export default routerSongCollaborator