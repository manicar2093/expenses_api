/*
  Warnings:

  - You are about to drop the column `day` on the `expenses` table. All the data in the column will be lost.
  - You are about to drop the column `month` on the `expenses` table. All the data in the column will be lost.
  - You are about to drop the column `year` on the `expenses` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "expenses" DROP COLUMN "day",
DROP COLUMN "month",
DROP COLUMN "year";
