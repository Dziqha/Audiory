import { PrismaClient } from "@prisma/client";
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
    log : [
        "query"
    ]
})


export async function AddPlaylistSong(req: any, res: any) {
    const { playlistId, songId } = req.body;
    if (!playlistId || !songId) {
        return res.status(400).json({ message: "Missing required fields" });
    }

    const playlist = await prisma.playlist.findUnique({
        where : {
            id : Number(playlistId)
        }
    });
    const song = await prisma.songs.findUnique({
        where : {
            id : Number(songId)
        }
    });
    if (!playlist || !song) {
        return res.status(404).json({ message: "Playlist or Song not found" });
    }

    try {
        const newPlaylistSong = await prisma.playlistSong.create({
            data : {
                playlistId : Number(playlistId),
                songId : Number(songId)
            }
        })

        return res.status(200).json({
            message : "PlaylistSong added successfully",
            data : newPlaylistSong
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error creating PlaylistSong"
        })
    }
}

export async function GetAllPlaylistsSongs(req : any, res : any) {
    const starttime = Date.now();
    const cachekey = process.env.CACHE_KEY_PLAYLISTSONG_ALL!;
    const client = await Initialize();
    const cachePlaylistSongs = await client.get(cachekey);
    if (cachePlaylistSongs) {
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message : "PlaylistSongs found successfully from cache",
            data : JSON.parse(cachePlaylistSongs),
            duration : duration
        })    
    }
    try {
        const findPlaylistSongs = await prisma.playlistSong.findMany();
        if (findPlaylistSongs.length === 0){
            return res.status(404).json({
                message : "PlaylistSongs not found"
            })
        }

        res.setHeader("X-Cache", "MISS");
        await client.setex(cachekey, 60, JSON.stringify(findPlaylistSongs));
        return res.status(200).json({
            message : "PlaylistSongs found successfully",
            data : findPlaylistSongs
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error finding PlaylistSongs"
        })
    }
}

export async function GetPlaylistSongById(req : any, res : any) {
    const starttime = Date.now();
    const cachekey = process.env.CACHE_KEY_PLAYLISTSONG_PREFIX!;
    const { id } = req.params;
    const client = await Initialize();
    const cachePlaylist = await client.get(cachekey + id);
    if (cachePlaylist) {
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message : "PlaylistSong found successfully from cache",
            data : JSON.parse(cachePlaylist),
            duration : duration
        })    
    }
    try {
        const findPlaylistSong = await prisma.playlistSong.findUnique({
            where : {
                id : Number(id)
            }
        })
        if (!findPlaylistSong) {
            return res.status(404).json({
                message : "PlaylistSong not found"
            })
        }

        res.setHeader("X-Cache", "MISS");
        await client.setex(cachekey + id, 60, JSON.stringify(findPlaylistSong));
        return res.status(200).json({
            message : "PlaylistSong found successfully",
            data : findPlaylistSong
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error finding PlaylistSong"
        })
    }
}

export async function UpdatePlaylistSong(req: any, res: any) {
    const { id } = req.params;
    const { playlistId, songId } = req.body;
    const dataPlaylistSong: any = {};

    if (playlistId !== undefined) {
        const findPlaylist = await prisma.playlist.findUnique({
            where: {
                id: Number(playlistId),
            },
        });
        if (!findPlaylist) {
            return res.status(404).json({
                message: "Playlist not found",
            });
        }
        dataPlaylistSong.playlistId = playlistId;
    }

    if (songId !== undefined) {
        const findSong = await prisma.songs.findUnique({
            where: {
                id: Number(songId),
            },
        });
        if (!findSong) {
            return res.status(404).json({
                message: "Song not found",
            });
        }
        dataPlaylistSong.songId = songId;
    }

    if (Object.keys(dataPlaylistSong).length === 0) {
        return res.status(400).json({
            message: "No valid fields to update",
        });
    }

    try {
        const updatedPlaylistSong = await prisma.playlistSong.update({
            where: {
                id: Number(id),
            },
            data: dataPlaylistSong,
        });
        return res.status(200).json({
            message: "PlaylistSong updated successfully",
            data: updatedPlaylistSong,
        });
    } catch (error) {
        console.error("Error updating PlaylistSong:", error);
        return res.status(500).json({
            message: "Error updating PlaylistSong",
        });
    }
}

export async function DeletePlaylistSong(req : any, res : any) {
    const { id } = req.params;
    const findPlaylistSong = await prisma.playlistSong.findUnique({
        where : {
            id : Number(id)
        }
    })
    if (!findPlaylistSong) {
        return res.status(404).json({
            message : "PlaylistSong not found"
        })
    }
    try {
        await prisma.playlistSong.delete({
            where : {
                id : Number(id)
            }
        })
        return res.status(200).json({
            message : "PlaylistSong deleted successfully",
        })
    }catch (error) {
        return res.status(500).json({
            message : "Error deleting PlaylistSong"
        })
    }
}