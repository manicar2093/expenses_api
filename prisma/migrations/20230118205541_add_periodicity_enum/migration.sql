-- CreateEnum
CREATE TYPE "Periodicity" AS ENUM ('Daily', 'Weekly', 'FourteenDaily', 'Paydaily', 'Monthly', 'BiMonthly', 'FourMonthly', 'SixMonthly', 'Yearly');

-- AlterTable
ALTER TABLE "recurrent_expenses" ADD COLUMN     "periodicity" "Periodicity" NOT NULL DEFAULT 'Monthly';
