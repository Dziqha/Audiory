import { PrismaClient } from "@prisma/client";
import * as grpc from '@grpc/grpc-js'
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
  log: ["query"],
});

export async function AddArtist(req: any, res: any) {
  const { name, debut_year } = req.body;
  if (!name || !debut_year) {
    return res.status(400).json({ message: "Missing required fields" });
  }

  const newArtist = await prisma.artist.create({
    data: {
      name,
      debutYear: debut_year,
    },
  });
  res.status(200).json({
    message: "Artist added successfully",
    data: newArtist,
  });
}

export async function GetALlArtist(req: any, res: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_ARTIST_ALL!;

  const client = await Initialize();
  const cacheArtist = await client.get(cachekey);
  if (cacheArtist) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Artists found successfully from cache",
      data: JSON.parse(cacheArtist),
      duration: duration,
    });
  }
  try {
    const artists = await prisma.artist.findMany();
    if (artists.length === 0) {
      return res.status(404).json({ message: "No artists found" });
    }
    await client.setex(cachekey, 60, JSON.stringify(artists));
    res.setHeader("X-Cache", "MISS");
    return res.status(200).json({
      message: "Artists found successfully",
      data: artists,
    });
  } catch (error) {
   return res.status(500).json({ message: "Internal server error" });
  }
}

export async function GetArtistById(req: any, res: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_ARTIST_PREFIX!;
  const { id } = req.params;
  const client = await Initialize();
  const cacheArtist = await client.get(cachekey + id);
  if (cacheArtist) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Artist found successfully from cache",
      data: JSON.parse(cacheArtist),
      duration: duration,
    });
  }
  try {
    const artist = await prisma.artist.findUniqueOrThrow({
      where: {
        id: Number(id),
      },
    });

    if (artist === null) {
      return res.status(404).json({ message: "Artist not found" });
    }
    await client.setex(cachekey + id, 60, JSON.stringify(artist));
    res.setHeader("X-Cache", "MISS");
    return res.status(200).json({
      message: "Artist found successfully",
      data: artist,
    });
  } catch (error) {
    res.status(404).json({ message: "Artist not found" });
  }
}

export async function UpdateArtist(req: any, res: any) {
  const { id } = req.params;
  const { name, debut_year } = req.body;
  try {
    const artist = await prisma.artist.updateMany({
      where: {
        id: Number(id),
      },
      data: {
        name,
        debutYear: debut_year,
      },
    });
    res.status(200).json({
      message: "Artist updated successfully",
      data: artist,
    });
  } catch (error) {
    res.status(404).json({ message: "Artist not found" });
  }
}

export async function DeleteArtist(req: any, res: any) {
  const { id } = req.params;
  try {
    const artist = await prisma.artist.delete({
      where: {
        id: Number(id),
      },
    });

    if (!artist) {
      return res.status(404).json({ message: "Artist not found" });
    }
    res.status(200).json({
      message: "Artist deleted successfully",
    });
  } catch (error) {
    return res.status(500).json({ message: "Internal server error" });
  }
}


export async function getArtistById(call: any, callback: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_ARTIST_PREFIX!;
  const { id } = call.request;
  const client = await Initialize();
  const cacheArtist = await client.get(cachekey + id);
  const responseMetadata = new grpc.Metadata();
  if (cacheArtist) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    responseMetadata.add("X-Cache", "HIT");
    const cacheData = JSON.parse(cacheArtist)
    console.log("Hit from cache", duration);
    return callback(null, {
      id: cacheData.id,
      name : cacheData.name,
      debutYear : cacheData.debutYear
  });
  }
  try {
      const artistId = call.request.id;

      if (!artistId) {
          return callback({
              code: grpc.status.INVALID_ARGUMENT,
              details: "Artist ID is required",
          });
      }

      const artist = await prisma.artist.findUnique({
          where: { id: Number(artistId) },
      });

      if (!artist) {
          return callback({
              code: grpc.status.NOT_FOUND,
              details: "artist not found",
          });
      }
      await client.setex(cachekey + id, 60, JSON.stringify(artist));
      responseMetadata.add("X-Cache", "MISS");

      return callback(null, {
          id: artist.id,
          name : artist.name,
          debutYear : artist.debutYear
      });

  } catch (error) {
      console.error('Error fetching Artist:', error);
      callback({
          code: grpc.status.INTERNAL,
          details: "Internal server error",
      });
  }
}

