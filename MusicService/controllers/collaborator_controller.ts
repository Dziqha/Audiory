import { PrismaClient } from "@prisma/client";
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
  log: ["query"],
});

export async function AddCollaborator(req: any, res: any) {
  const { name, roleType } = req.body;
  const validateRoleType = [
    "Producer",
    "Composer",
    "Featured_Artist",
    "Lyricist",
    "Arranger",
    "Instrumentalist",
    "Engineer",
  ];
  if (!name || !roleType)
    return res.status(404).json({ message: "Missing required fields" });
  if (!validateRoleType.includes(roleType))
    return res.status(400).json({ message: "Invalid role type" });
  try {
    const newCollaborator = await prisma.collaborator.create({
      data: {
        name,
        roleType,
      },
    });

    return res.status(200).json({
      message: "Collaborator added successfully",
      data: newCollaborator,
    });
  } catch (error) {
    return res.status(500).json({
      message: "Collaborator not added",
    });
  }
}


export async function GetAllCollaborators(req : any, res : any){
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_COLLABORATOR_ALL!;
  const client = await Initialize();
  const cacheCollaborator = await client.get(cachekey);
  if (cacheCollaborator) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Collaborators found successfully from cache",
      data: JSON.parse(cacheCollaborator),
      duration: duration,
    });
  }
  try {
    const findCollaborator =  await prisma.collaborator.findMany()

    if (findCollaborator.length === 0) {
      return res.status(404).json({
        message : "No collaborators found"
      })
    }

    res.setHeader("X-Cache", "MISS");
    await client.setex(cachekey, 60, JSON.stringify(findCollaborator));
    return res.status(200).json({
      message : "Collaborators found successfully",
      data : findCollaborator
    })
  }catch (error) {
    return res.status(500).json({
      message : "Error finding collaborators"
    })
  }
}

export async function GetCollaboratorById(req : any, res : any){
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_PLAYLISTSONG_PREFIX!;
  const { id } = req.params;
  const client = await Initialize();
  const cacheCollaborator = await client.get(cachekey + id);
  if (cacheCollaborator) {
    const endtime = Date.now();
    const duration = endtime - starttime;    
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Collaborator found successfully from cache",
      data: JSON.parse(cacheCollaborator),
      duration: duration,
    });
  }
  try {
    const findCollaborator = await prisma.collaborator.findUnique({
      where : {
        id : Number(id)
      }
    })

    if (!findCollaborator) {
      return res.status(404).json({
        message : "Collaborator not found"
      })
    }
    res.setHeader("X-Cache", "MISS");
    await client.setex(cachekey + id, 60, JSON.stringify(findCollaborator));

    return res.status(200).json({
      message : "Collaborator found successfully",
      data : findCollaborator
    })
  }catch(error) {
    return res.status(500).json({
      message : "Error finding collaborator"
    })
  }
}

export async function UpdateCollaborator(req : any, res : any){
  const { id } = req.params;
  const { name, roleType } = req.body;

  const findCollaborator = await prisma.collaborator.findUnique({
    where : {
      id :  Number(id)
    }
  })

  if(!findCollaborator) {
    return res.status(404).json({
      message : "Collaborator not found"
    })
  }

  const dataCollaborator : any = {};
  if (name !== undefined) dataCollaborator.name = name;
  if (roleType !== undefined) dataCollaborator.roleType = roleType;

  try {
    const updateCollaborator =  await prisma.collaborator.updateMany({
      where : {
        id : Number(id)
      },
      data : dataCollaborator
    })
    return res.status(200).json({
      message : "Collaborator updated successfully",
      data : updateCollaborator
    })
  }catch (error) {
    return res.status(500).json({
      message : "Error updating collaborator"
    })
  }
}


export async function DeleteCollaborator(req : any, res : any) {
  const { id } =  req.params;

  const findCollaborator = await prisma.collaborator.findUnique({
    where : {
      id : Number(id)
    }
  })

  if (!findCollaborator) {
    return res.status(404).json({
      message : "Collaborator not found"
    })
  }

  try {
    await prisma.collaborator.delete({
      where : {
        id : Number(id)
      }
    })
    return res.status(200).json({
      message : "Collaborator deleted successfully",
    })
  }catch (error) {
    return res.status(500).json({
      message : "Error deleting collaborator"
    })
  }
}
