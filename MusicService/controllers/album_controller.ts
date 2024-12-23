import { PrismaClient } from "@prisma/client";
import * as grpc from '@grpc/grpc-js'
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
  log: ["query"],
});

export async function AddAlbums(req: any, res: any) {
  const { title, releaseYear, artistId, genreId } = req.body;

  if (!title || !releaseYear || !artistId || !genreId) {
    return res.status(400).json({ message: "Missing required fields" });
  }

  const artist = await prisma.artist.findUnique({
    where: {
      id: Number(artistId),
    },
  });

  const genre = await prisma.genre.findUnique({
    where: {
      id: Number(genreId),
    },
  });

  if (!artist || !genre) {
    return res.status(404).json({ message: "Artist or Genre not found" });
  }

  try {
    const newAlbum = await prisma.album.create({
      data: {
        title,
        releaseYear,
        artistId: Number(artistId),
        genreId: Number(genreId),
      },
    });

    return res.status(200).json({
      message: "Album created successfully",
      data: newAlbum,
    });
  } catch (error) {
    return res.status(500).json({
      message: "Internal server error",
    });
  }
}

export async function GetAllAlbum(req: any, res: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_ALBUM_ALL!;

  const client = await Initialize();
  const cacheAlbum = await client.get(cachekey);

  if (cacheAlbum) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Albums found successfully from cache",
      data: JSON.parse(cacheAlbum),
      duration: duration,
    });
  }
  try {
    const albums = await prisma.album.findMany();
    if (albums.length === 0) {
      return res.status(404).json({
        message: "No albums found",
      });
    }

    res.setHeader("X-Cache", "MISS");
    await client.setex(cachekey, 60, JSON.stringify(albums));
    return res.status(200).json({
      message: "Albums successfully retrieved",
      data: albums,
    });
  } catch (error) {
    return res.status(500).json({
      message: "Internal server error",
    });
  }
}

export async function GetAlbumById(req: any, res: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_ALBUM_PREFIX!;
  const { id } = req.params;

  const client = await Initialize();
  const cacheAlbum = await client.get(cachekey + id);

  if (cacheAlbum) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Album found successfully from cache",
      data: JSON.parse(cacheAlbum),
      duration: duration,
    });
  }
  try {
    const album = await prisma.album.findUnique({
      where: {
        id: Number(id),
      },
    });
    res.setHeader("X-Cache", "MISS");
    await client.setex(cachekey + id, 60, JSON.stringify(album));
    return res.status(200).json({
      message: "Album successfully retrieved",
      data: album,
    });
  } catch (error) {
    return res.status(500).json({
      message: "Internal server error",
    });
  }
}

export async function UpdateAlbum(req: any, res: any) {
    const { id } = req.params;
    const { title, releaseYear, artistId, genreId } = req.body;

    const updateData: any = {};
    if (title !== undefined) updateData.title = title;
    if (releaseYear !== undefined) updateData.releaseYear = releaseYear;
    if (artistId !== undefined) updateData.artistId = Number(artistId);
    if (genreId !== undefined) updateData.genreId = Number(genreId);

    try {
        const existingAlbum = await prisma.album.findUnique({
            where: { id: Number(id) },
        });

        if (!existingAlbum) {
            return res.status(404).json({ message: "Album not found" });
        }

        const updatedAlbum = await prisma.album.update({
            where: { id: Number(id) },
            data: updateData,
        });

        return res.status(200).json({
            message: "Album updated successfully",
            data: updatedAlbum,
        });
    } catch (error) {
        return res.status(500).json({
            message: "Error updating album",
        });
    }
}


export async function DeleteAlbum(req : any, res : any) {
    const { id } = req.params;
    try {
        await prisma.album.delete({
            where: {
                id: Number(id),
            },
        });
        res.status(200).json({
            message: "Album deleted successfully",
        });
    } catch (error) {
        res.status(500).json({
            message: "Internal server error",
        });
    }
}


export async function getAlbumById (call: any, callback: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_ALBUM_PREFIX!;
  const { id } = call.request;

  const client = await Initialize();
  const cacheAlbum = await client.get(cachekey + id);
  const responseMetadata = new grpc.Metadata();

  if (cacheAlbum) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    responseMetadata.add("X-Cache", "HIT");
    const cacheData = JSON.parse(cacheAlbum)
    console.log("Hit from cache", duration);
    callback(null,{
      id : cacheData.id,
      title : cacheData.title,
      releaseYear : cacheData.releaseYear,
      artistId : cacheData.artistId,
      genreId : cacheData.genreId
    })
  }
  try {
    const albumId = call.request.id

    if (!albumId) {
      return callback({
        code : grpc.status.INVALID_ARGUMENT,
        details : "Album ID is required"
      })
    }

    const album = await prisma.album.findUnique({
      where : {
        id : Number(albumId)
      }
    })

    if (!album) {
      return callback({
        code : grpc.status.NOT_FOUND,
        details : "Album not found"
      })
    }

    responseMetadata.add("X-Cache", "MISS");
    await client.setex(cachekey + id, 60, JSON.stringify(album));
    callback(null, {
      id : album.id,
      title : album.title,
      releaseYear : album.releaseYear,
      artistId : album.artistId,
      genreId : album.genreId
    })
  }catch (error) {
    console.log("Failed to fetch album:", error);
    return callback({
      code : grpc.status.INTERNAL,
      details : "Failed to fetch album"
    })
  }
}
