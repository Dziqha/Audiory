import { PrismaClient } from "@prisma/client";
import * as grpc from '@grpc/grpc-js'
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
    log : [
        "query"
    ]
})

export async function AddSongs(req : any, res: any) {
    const {title, releaseDate, duration, albumId, artistId, genreId, isPublished} = req.body;

    if (!title || !releaseDate || !duration || !albumId || !artistId || !genreId || !isPublished) {
        return res.status(400).json({
            message : "Missing required fields"
        })
    }
    try {
        const newSong = await prisma.songs.create({
            data: {
                title,
                releaseDate,
                duration,
                albumId: Number(albumId),
                artistId: Number(artistId),
                genreId: Number(genreId),
                isPublished
            }
        })

        return res.status(200).json({
            message : "Song added successfully",    
            data : newSong
        })
    }catch(error){
        return res.status(500).json({
            message : "Internal server error"
        })
    }
}

export async function GetAllSongs(req : any, res: any) {
    const starttime = Date.now();
    const cachekey = process.env.CACHE_KEY_SONG_ALL!;

    const client = await Initialize();
    const cachesong = await client.get(cachekey);
    if (cachesong) {
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message : "Songs successfully retrieved from cache",
            data : JSON.parse(cachesong),
            duration : duration
        })
    }
    try {
        const songs = await prisma.songs.findMany();
        if (songs.length === 0) {
            return res.status(404).json({
                message : "No songs found"
            })
        }
        await client.setex(cachekey, 60, JSON.stringify(songs));
        res.setHeader("X-Cache", "MISS");
        return res.status(200).json({
            message : "Songs successfully retrieved",
            data : songs
        })
    }catch (error){
        return res.status(500).json({
            message : "Internal server error"
        })
    }
}

export async function GetSongById(req : any, res: any){
    const starttime = Date.now();
    const cachekey = process.env.CACHE_KEY_SONG_PREFIX!;
    const { id } = req.params;

    const client = await Initialize();
    const cachesong = await client.get(cachekey + id);
    if (cachesong) {
        const endtime = Date.now();
        const duration = endtime - starttime;
        res.setHeader("X-Cache", "HIT");
        console.log("Hit from cache", duration);
        return res.status(200).json({
            message : "Song successfully retrieved from cache",
            data : JSON.parse(cachesong),
            duration : duration
        })
    }
    try{
        const song = await prisma.songs.findUnique({
            where : {
                id : Number(id),
            }
        })

        if (!song) {
            return res.status(404).json({
                message : "Song not found"
            })
        }

        await client.setex(cachekey + id, 60, JSON.stringify(song));
        res.setHeader("X-Cache", "MISS");

        return res.status(200).json({
            message : "Song successfully retrieved",
            data : song
        })
    }catch (error) {
        return res.status(500).json({
            message : "Internal server error"
        })
    }
}

export async function UpdateSong(req : any, res : any) {
    const { id } = req.params;
    const {title, releaseDate, duration, albumId, artistId, genreId, isPublished} = req.body;

    const dataSong: any = {};
    if (title !== undefined) dataSong.title =  title;
    if (releaseDate !== undefined) dataSong.releaseDate =  releaseDate;
    if (duration !== undefined) dataSong.duration =  duration;
    if (albumId !== undefined)  dataSong.albumId = Number(albumId);
    if (artistId !== undefined) dataSong.artistId = Number(artistId);
    if (genreId !== undefined) dataSong.genreId = Number(genreId);
    if (isPublished !== undefined) dataSong.isPublished = isPublished;

    try {
        const findSong =  await prisma.songs.findUnique({
            where : {
                id : Number(id)
            }
        })

        if (!findSong) {
            return res.status(404).json({
                message : "Song not found"
            })
        }

        const updateSong =  await prisma.songs.updateMany({
            where : {
                id : Number(id)
            },
            data : dataSong,
        })

        return res.status(200).json({
            message : "Song updated successfully",
            data : updateSong
        })
    }catch (error) {
        return res.status(500).json({
            message : "Internal server error"
        })
    }
}

export async function DeleteSong(req : any, res : any) {
    const { id } = req.params;
    try {
        const findSong =  await prisma.songs.findUnique({
            where : {
                id : Number(id)
            }
        })

        if (!findSong) {
            return res.status(404).json({
                message : "Song not found"
            })
        }
        
        await prisma.songs.delete({
            where : {
                id : Number(id)
            }
        })

        return res.status(200).json({
            message : "Song deleted successfully",
        })
    }catch (error) {
        return res.status(500).json({
            message : "Internal server error"
        })
    }
}


export async function getSongById(call: any, callback: any) {
    try {
        const songId = call.request.id;

        if (!songId) {
            return callback({
                code: grpc.status.INVALID_ARGUMENT,
                details: "Song ID is required",
            });
        }

        const song = await prisma.songs.findUnique({
            where: { id: Number(songId) },
        });

        if (!song) {
            return callback({
                code: grpc.status.NOT_FOUND,
                details: "Song not found",
            });
        }

        const genre = song.genreId
            ? await prisma.genre.findUnique({
                where: { id: song.genreId },
            })
            : null;
        
        const artist = song.artistId
            ? await prisma.artist.findUnique({
                where: { id: song.artistId },
            })
            : null;

        const album = song.albumId
            ? await prisma.album.findUnique({
                where: { id: song.albumId },
            })
            : null;

        callback(null, {
            id: song.id,
            title: song.title,
            artistId: artist ? artist.id : null,
            albumId: album ? album.id : null,
            genreId: genre ? genre.id : null,
            duration: song.duration,
        });

    } catch (error) {
        console.error('Error fetching song:', error);
        callback({
            code: grpc.status.INTERNAL,
            details: "Internal server error",
        });
    }
}
