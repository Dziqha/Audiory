import { PrismaClient } from "@prisma/client";
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
    log : [
        "query"
    ]
})


export async function AddPlaylist(req : any, res : any) {
    const { name, createdAt, userId} = req.body;
    if (!name || !createdAt || !userId) {
        return res.status(400).json({ message: "Missing required fields" });
    }
    try {
        const newPlaylist =  await prisma.playlist.create({
            data : {
                name,
                createdAt,
                userId : Number(userId)
            }
        })
        return res.status(200).json({
            message : "Playlist added successfully",
            data : newPlaylist
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error creating playlist"
        })
    }
}

export async function GetAllPlaylists(req : any, res : any) {
    const starttime = Date.now();
    const cachekey =  process.env.CACHE_KEY_PLAYLIST_ALL!;
    const client = await Initialize();
    const cachePlaylist = await client.get(cachekey);
    if (cachePlaylist) {
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message : "Playlists found successfully from cache",
            data : JSON.parse(cachePlaylist),
            duration : duration
        })    
    }
    try {
        const findPlaylists = await prisma.playlist.findMany();
        if (findPlaylists.length === 0) {
            return res.status(404).json({
                message : "No playlists found"
            })
        }

        res.setHeader("X-Cache", "MISS");
        await client.setex(cachekey, 60, JSON.stringify(findPlaylists));
        return res.status(200).json({
            message : "Playlists found successfully",
            data : findPlaylists
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error finding playlists"
        })
    }
}

export async function GetPlaylistsById(req : any, res : any) {
    const starttime = Date.now();
    const cachekey = process.env.CACHE_KEY_PLAYLIST_PREFIX!;
    const { id } = req.params;

    const client =  await Initialize();
    const cachePlaylist = await client.get(cachekey + id);
    if (cachePlaylist) {
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message : "Playlist found successfully from cache",
            data : JSON.parse(cachePlaylist),
            duration : duration
        })    
    }

    try {
        const findPlaylist = await prisma.playlist.findUnique({
            where : {
                id : Number(id)
            }
        })

        if (!findPlaylist) {
            return res.status(404).json({
                message : "Playlist not found"
            })
        }

        res.setHeader("X-Cache", "MISS");
        await client.setex(cachekey + id, 60, JSON.stringify(findPlaylist));

        return res.status(200).json({
            message : "Playlist found successfully",
            data : findPlaylist
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error finding playlist"
        })
    }
}


export async function UpdatePlaylist(req : any, res : any) {
    const { id } = req.params;
    const { name, createdAt, userId} = req.body;

    const findPlaylist = await prisma.playlist.findUnique({
        where : {
            id : Number(id)
        }
    })

    if (!findPlaylist) {
        return res.status(404).json({
            message : "Playlist not found"
        })
    }

    const dataPlaylist : any = {};
    if (name !== undefined) dataPlaylist.name = name;
    if (createdAt !== undefined) dataPlaylist.createdAt = createdAt;
    if (userId !== undefined) dataPlaylist.userId = Number(userId);

    try {
        const updatePlaylist = await prisma.playlist.updateMany({
            where : {
                id :  Number(id)
            },

            data : dataPlaylist
        })
        return res.status(200).json({
            message : "Playlist updated successfully",
            data : updatePlaylist
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error updating playlist"
        })
    }
}

export async function DeletePlaylist(req : any, res : any ){
    const { id } = req.params;

    const findPlaylist = await prisma.playlist.findUnique({
        where : {
            id : Number(id)
        }
    })

    if (!findPlaylist) {
        return res.status(404).json({
            message : "Playlist not found"
        })
    }

    try {
        await prisma.playlist.delete({
            where : {
                id : Number(id)
            }
        })

        return res.status(200).json({
            message : "Playlist deleted successfully"
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error deleting playlist"
        })
    }
}