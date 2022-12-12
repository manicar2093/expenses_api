-- AlterTable
ALTER TABLE "expenses" ADD COLUMN     "user_id" UUID;

-- AlterTable
ALTER TABLE "incomes" ADD COLUMN     "user_id" UUID;

-- AlterTable
ALTER TABLE "recurrent_expenses" ADD COLUMN     "user_id" UUID;

-- AddForeignKey
ALTER TABLE "expenses" ADD CONSTRAINT "expenses_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "incomes" ADD CONSTRAINT "incomes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "recurrent_expenses" ADD CONSTRAINT "recurrent_expenses_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;
