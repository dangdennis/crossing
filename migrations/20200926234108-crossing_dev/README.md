# Migration `20200926234108-crossing_dev`

This migration has been generated by Dennis Dang at 9/26/2020, 11:41:08 PM.
You can check out the [state of the schema](./schema.prisma) after the migration.

## Database Steps

```sql
CREATE UNIQUE INDEX "Action.eventId_position" ON "public"."Action"("eventId","position")

CREATE UNIQUE INDEX "AvatarsOnRaids.raidId_avatarId" ON "public"."AvatarsOnRaids"("raidId","avatarId")

CREATE UNIQUE INDEX "Event.storyId_sequence" ON "public"."Event"("storyId","sequence")

CREATE UNIQUE INDEX "MessagesOnActions.messageId_actionId" ON "public"."MessagesOnActions"("messageId","actionId")

CREATE UNIQUE INDEX "RaidBossesOnRaids.raidId_raidBossId" ON "public"."RaidBossesOnRaids"("raidId","raidBossId")

CREATE UNIQUE INDEX "StoriesOnRaids.raidId_storyId" ON "public"."StoriesOnRaids"("raidId","storyId")

ALTER INDEX "public"."Avatar_userId" RENAME TO "Avatar.userId"
```

## Changes

```diff
diff --git schema.prisma schema.prisma
migration 20200926232750-crossing_dev..20200926234108-crossing_dev
--- datamodel.dml
+++ datamodel.dml
@@ -1,7 +1,7 @@
 datasource db {
     provider = "postgresql"
-    url = "***"
+    url      = "postgresql://postgres:postgres@localhost:5432/crossing_dev"
 }
 generator db {
     provider = "go run github.com/prisma/prisma-client-go"
@@ -25,9 +25,9 @@
     createdAt      DateTime         @default(now())
     deletedAt      DateTime?
     updatedAt      DateTime         @updatedAt
     User           User             @relation(fields: [userId], references: [id])
-    userId         Int
+    userId         Int              @unique
     AvatarsOnRaids AvatarsOnRaids[]
 }
 model Raid {
@@ -66,8 +66,9 @@
     avatar    Avatar    @relation(fields: [avatarId], references: [id])
     avatarId  Int
     @@id([raidId, avatarId])
+    @@unique([raidId, avatarId])
 }
 model RaidBossesOnRaids {
     createdAt  DateTime  @default(now())
@@ -78,8 +79,9 @@
     raidBoss   RaidBoss  @relation(fields: [raidBossId], references: [id])
     raidBossId Int
     @@id([raidId, raidBossId])
+    @@unique([raidId, raidBossId])
 }
 // A Story is a series of events that can occur for a given raid
 model Story {
@@ -101,8 +103,9 @@
     story     Story     @relation(fields: [storyId], references: [id])
     storyId   Int
     @@id([raidId, storyId])
+    @@unique([raidId, storyId])
 }
 // An Event represents a moment in a story in which Avatars can fulfill Actions
 model Event {
@@ -114,8 +117,11 @@
     storyId   Int
     // The sequence dictates the order of events to occur in a story
     sequence  Int
     Action    Action[]
+
+
+    @@unique([storyId, sequence])
 }
 // An Action prescribes a message and the raid member position required to fulfill the action.
 // This allows for customized actions per raid member. Wowzah!
@@ -127,8 +133,10 @@
     position          Int
     event             Event               @relation(fields: [eventId], references: [id])
     eventId           Int
     MessagesOnActions MessagesOnActions[]
+
+    @@unique([eventId, position])
 }
 // A Message holds message contents.
 model Message {
@@ -149,5 +157,6 @@
     action    Action    @relation(fields: [actionId], references: [id])
     actionId  Int
     @@id([messageId, actionId])
+    @@unique([messageId, actionId])
 }
```

