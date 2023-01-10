/*
  Warnings:

  - You are about to drop the column `day` on the `incomes` table. All the data in the column will be lost.
  - You are about to drop the column `month` on the `incomes` table. All the data in the column will be lost.
  - You are about to drop the column `year` on the `incomes` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "incomes" DROP COLUMN "day",
DROP COLUMN "month",
DROP COLUMN "year";
