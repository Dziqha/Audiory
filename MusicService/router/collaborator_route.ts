import { AddCollaborator } from "../controllers/collaborator_controller";
import { GetAllCollaborators } from "../controllers/collaborator_controller";
import { GetCollaboratorById } from "../controllers/collaborator_controller";
import { UpdateCollaborator } from "../controllers/collaborator_controller";
import { DeleteCollaborator } from "../controllers/collaborator_controller";
import { AuthenticateAdmin } from "../middleware/auth";
import { Router } from "express";
const routerCollaborator = Router();


routerCollaborator.post('/add', AuthenticateAdmin, AddCollaborator)
routerCollaborator.get('/all', GetAllCollaborators)
routerCollaborator.get('/get/:id', GetCollaboratorById)
routerCollaborator.put('/update/:id', AuthenticateAdmin, UpdateCollaborator)
routerCollaborator.delete('/delete/:id', AuthenticateAdmin, DeleteCollaborator)

export default routerCollaborator