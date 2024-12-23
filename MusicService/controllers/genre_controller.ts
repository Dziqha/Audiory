import { PrismaClient } from "@prisma/client";
import dotenv from "dotenv";
import { Initialize } from "../configs/redis";

dotenv.config()
const prisma = new PrismaClient({
  log: ["query"],
});

export async function AddGenre(req: any, res: any) {
  const { name, description } = req.body;
  if (!name || !description) {
    return res.status(400).json({ message: "Missing required fields" });
  }

  try {
    const newGenre = await prisma.genre.create({
      data: {
        name,
        description,
      },
    });
    return res.status(200).json({
      message: "Genre added successfully",
      data: newGenre,
    });
  } catch (error) {
    return res.status(500).json({ message: "Internal server error" });
  }
}

export async function GetALlGenre(req: any, res: any) {
  const starttime = Date.now();
  const cachekey = process.env.CACHE_KEY_GENRE_ALL!;

  const client = await Initialize();
  const cacheGenre = await client.get(cachekey);
  if (cacheGenre) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Genres found successfully from cache",
      data: JSON.parse(cacheGenre),
      duration: duration,
    });
  }
  try {
    const genres = await prisma.genre.findMany();
    if (genres.length === 0) {
      return res.status(404).json({ message: "No genres found" });
    }
    res.setHeader("X-Cache", "MISS");
    await client.setex(cachekey, 60, JSON.stringify(genres));
    return res.status(200).json({
      message: "Genres found successfully",
      data: genres,
    });
  } catch (error) {
    res.status(500).json({ message: "Internal server error" });
  }
}

export async function GetGenreById(req: any, res: any) {
  const starttime =  Date.now();
  const cachekey = process.env.CACHE_KEY_GENRE_PREFIX!;
  const { id } = req.params;

  const client = await Initialize();
  const cacheGenre = await client.get(cachekey + id);
  if (cacheGenre) {
    const endtime = Date.now();
    const duration = endtime - starttime;
    res.setHeader("X-Cache", "HIT");
    console.log("Hit from cache", duration);
    return res.status(200).json({
      message: "Genre found successfully from cache",
      data: JSON.parse(cacheGenre),
      duration: duration,
    });
  }
  try {
    const genre = await prisma.genre.findUniqueOrThrow({
      where: {
        id: Number(id),
      },
    });
    if (genre === null) {
      return res.status(404).json({ message: "Genre not found" });
    }
    res.setHeader("X-Cache", "MISS");
    await client.setex(cachekey + id, 60, JSON.stringify(genre));
    return res.status(200).json({
      message: "Genre found successfully",
      data: genre,
    });
  } catch (error) {
    res.status(404).json({ message: "Genre not found" });
  }
}

export async function UpdateGenre(req: any, res: any) {
  const { id } = req.params;
  const { name, description } = req.body;
  try {
    const genre = await prisma.genre.updateMany({
      where: {
        id: Number(id),
      },
      data: {
        name,
        description,
      },
    });
    if (genre.count === 0) {
      return res.status(404).json({ message: "Genre not found" });
    }
    res.status(200).json({
      message: "Genre updated successfully",
      data: genre,
    });
  } catch (error) {
    res.status(500).json({ message: "Internal server error" });
  }
}

export async function DeleteGenre(req: any, res: any) {
  const { id } = req.params;
  try {
    const genre = await prisma.genre.delete({
      where: {
        id: Number(id),
      },
    });
    res.status(200).json({
      message: "Genre deleted successfully",
    });
  } catch (error) {
    return res.status(404).json({ message: "Genre not found" });
  }
}
