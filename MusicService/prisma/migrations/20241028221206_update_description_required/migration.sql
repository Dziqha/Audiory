/*
  Warnings:

  - Made the column `description` on table `genre` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE `genre` MODIFY `description` VARCHAR(191) NOT NULL;
