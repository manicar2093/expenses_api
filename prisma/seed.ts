import {PrismaClient} from '@prisma/client'

const prisma = new PrismaClient()

async function updateCatalog(catalog:Array<any>, entity: any) {
    catalog.forEach(async i => {
        await entity.upsert({
            where: {
                id: i.id,
            },
            update: {
                name: i.name
            },
            create: i,
        });
    });
}

async function main() {
    
}

main()
    .catch(e => {
        throw e
    })
    .finally(async () => {
        await prisma.$disconnect()
    });
