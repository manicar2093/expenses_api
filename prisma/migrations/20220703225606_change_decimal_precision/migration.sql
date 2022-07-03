/*
  Warnings:

  - You are about to alter the column `amount` on the `Expenses` table. The data in that column could be lost. The data in that column will be cast from `Decimal(65,30)` to `Decimal(9,2)`.
  - You are about to alter the column `amount` on the `Incomes` table. The data in that column could be lost. The data in that column will be cast from `Decimal(65,30)` to `Decimal(9,2)`.

*/
-- AlterTable
ALTER TABLE "Expenses" ALTER COLUMN "amount" SET DATA TYPE DECIMAL(9,2);

-- AlterTable
ALTER TABLE "Incomes" ALTER COLUMN "amount" SET DATA TYPE DECIMAL(9,2);
