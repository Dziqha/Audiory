import { PrismaClient } from "@prisma/client";
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
    log : [
        "query"
    ]
})


export async function AddSongCollaborator(req: any, res: any){
    const { songId, collaboratorId } = req.body;

    if (!songId || !collaboratorId) {
        return res.status(404).json({
            message: "Missing required fields"
        })
    }

    const findSong = await prisma.songs.findUnique({
        where : {
            id : Number(songId)
        }
    })

    const findCollaborator = await prisma.collaborator.findUnique({
        where : {
            id : Number(collaboratorId)
        }
    })

    if (!findSong || !findCollaborator) {
        return res.status(404).json({
            message : "Song or Collaborator not found"
        })
    }

    try {
        const newSongCollaborator = await prisma.songCollaborator.create({
            data : {
                songId : Number(songId),
                CollaboratorId : Number(collaboratorId)
            }
        })

        return res.status(200).json({
            message : "SongCollaborator created successfully",
            data : newSongCollaborator
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error creating SongCollaborator"
        })
    }
}

export async function FindAllSongCollaborator(req : any, res : any) {
    const starttime = Date.now();
    const cachekey = process.env.CACHE_KEY_SONGCOLLABORATOR_ALL!;
    const client = await Initialize();
    const cacheSongCollaborator = await client.get(cachekey);
    if (cacheSongCollaborator){
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message: "SongCollaborator found successfully from cache",
            data: JSON.parse(cacheSongCollaborator),
            duration: duration,
        });
    }
    try {
        const findSongCollaborator = await prisma.songCollaborator.findMany();

        if (findSongCollaborator.length === 0) return res.status(404).json({message : "SongCollaborator not found"});

        res.setHeader("X-Cache", "MISS");
        await client.setex(cachekey, 60, JSON.stringify(findSongCollaborator));
        return res.status(200).json({
            message : "SongCollaborator found successfully",
            data : findSongCollaborator
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error finding SongCollaborator"
        })
    }
}

export async function FindSongCollaboratorById(req : any, res : any) {
    const starttime = Date.now();
    const cachekey =  process.env.CACHE_KEY_SONGCOLLABORATOR_PREFIX!;
    const { id } = req.params;
    const client = await Initialize();
    const cacheSongCollaborator = await client.get(cachekey + id);
    if (cacheSongCollaborator){
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message: "SongCollaborator found successfully from cache",
            data: JSON.parse(cacheSongCollaborator),
            duration: duration,
        });
    }
    try {
        const findSongCollaborator = await prisma.songCollaborator.findUnique({
            where : {
                id : Number(id)
            }
        })

        if (!findSongCollaborator) return res.status(404).json({message : "SongCollaborator not found"});

        res.setHeader("X-Cache", "MISS");
        await client.setex(cachekey + id, 60, JSON.stringify(findSongCollaborator));
        return res.status(200).json({
            message : "SongCollaborator found successfully",
            data : findSongCollaborator
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error finding SongCollaborator"
        })
    }
}

export async function UpdateSongCollaborator(req : any, res : any) {
    const { id } = req.params;
    const { songId, collaboratorId } = req.body;
    const dataSongCollaborator : any = {};

    if (songId !== undefined){
        const song = await prisma.songs.findUnique({
            where : {
                id : Number(songId)
            }
        })
        if (!song) {
            return res.status(404).json({
                message : "Song not found"
            })
        }
        dataSongCollaborator.songId = songId;
    }

    if (collaboratorId !== undefined) {
        const collaborator = await prisma.collaborator.findUnique({
            where : {
                id : Number(collaboratorId)
            }
        })

        if (!collaborator){
            return res.status(404).json({
                message : "Collaborator not found"
            })
        }
        dataSongCollaborator.CollaboratorId = collaboratorId;
    }

    try {
        const findSongCollaborator = await prisma.songCollaborator.findUnique({
            where : {
                id : Number(id)
            }
        })

        if (!findSongCollaborator) return res.status(404).json({message : "SongCollaborator not found"});

        const updateSongCollaborator = await prisma.songCollaborator.update({
            where : {
                id : Number(id)
            },
            data : dataSongCollaborator
        })

        return res.status(200).json({
            message : "SongCollaborator updated successfully",
            data : updateSongCollaborator
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error updating SongCollaborator"
        })
    }
}


export async function DeleteSongCollaborator(req : any, res : any) {
    const { id } = req.params;
    try {
        const findSongCollaborator = await prisma.songCollaborator.findUnique({
            where : {
                id : Number(id)
            }
        })

        if (!findSongCollaborator) return res.status(404).json({message : "SongCollaborator not found"});

        await prisma.songCollaborator.delete({
            where : {
                id : Number(id)
            }
        })

        return res.status(200).json({
            message : "SongCollaborator deleted successfully"
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error deleting SongCollaborator"
        })
    }
}