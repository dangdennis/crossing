package seeder

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v5"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/repositories/messages"
	"github.com/dangdennis/crossing/repositories/users"
)

// Run runs the seeder
func Run() {
	db := prisma.Client()
	random, _ := rand.Prime(rand.Reader, 128)
	ctx := context.Background()

	// Create a couple users
	gofakeit.Seed(random.Int64())
	user1, err := users.CreateUser(db, users.UserAttrs{
		DiscordUserID:   strconv.FormatUint(uint64(gofakeit.Number(10000000, 90000000)), 10),
		Email:           toPtrString(gofakeit.Email()),
		DiscordUsername: toPtrString(gofakeit.Username()),
		FirstName:       toPtrString(gofakeit.FirstName()),
		LastName:        toPtrString(gofakeit.LastName()),
	})
	handleError(err)

	user2, err := users.CreateUser(db, users.UserAttrs{
		DiscordUserID:   strconv.FormatUint(uint64(gofakeit.Number(10000000, 90000000)), 10),
		Email:           toPtrString(gofakeit.Email()),
		DiscordUsername: toPtrString(gofakeit.Username()),
		FirstName:       toPtrString(gofakeit.FirstName()),
		LastName:        toPtrString(gofakeit.LastName()),
	})
	handleError(err)

	// Create the avatars for the new users
	avatar1, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(user1.ID),
		),
	).Exec(ctx)
	handleError(err)

	avatar2, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(user2.ID),
		),
	).Exec(ctx)
	handleError(err)

	// Create a couple starter raid bosses
	bossLichKing, err := db.RaidBoss.CreateOne(
		prisma.RaidBoss.Name.Set("Arthas Menethil, The Lich King"),
		prisma.RaidBoss.Image.Set("https://cdn.vox-cdn.com/thumbor/k6m7tw54mdYa2yJoYbk3FuIYFZg=/0x0:1024x576/1920x0/filters:focal(0x0:1024x576):format(webp):no_upscale()/cdn.vox-cdn.com/uploads/chorus_asset/file/19748343/155054_the_lich_king.jpg"),
	).Exec(ctx)
	handleError(err)

	bossAlienQueen, err := db.RaidBoss.CreateOne(
		prisma.RaidBoss.Name.Set("The Alien Queen"),
		prisma.RaidBoss.Image.Set("https://vignette.wikia.nocookie.net/avp/images/7/74/Promo07.PNG/revision/latest?cb=20120527102557"),
	).Exec(ctx)
	handleError(err)

	now := time.Now()
	later := time.Now().Add(49 * time.Hour)

	// Create a couple raids and story
	// The first is inactive. The second is active, open for another 7 days from time of seed.
	story1, err := db.Story.CreateOne().Exec(ctx)
	handleError(err)

	raid1, err := db.Raid.CreateOne(
		prisma.Raid.Story.Link(
			prisma.Story.ID.Equals(story1.ID),
		),
		prisma.Raid.Active.Set(false),
		prisma.Raid.PlayerLimit.Set(20),
		prisma.Raid.StartTime.Set(now.Add(1*time.Second)),
		prisma.Raid.EndTime.Set(later),
		prisma.Raid.CompletionProgress.Set(1.0),
	).Exec(ctx)
	handleError(err)

	// Create the Alien Queen story
	alienStory, err := db.Story.CreateOne().Exec(ctx)
	handleError(err)

	raid2, err := db.Raid.CreateOne(
		prisma.Raid.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID),
		),
		prisma.Raid.Active.Set(true),
		prisma.Raid.PlayerLimit.Set(20),
		prisma.Raid.StartTime.Set(now),
		prisma.Raid.EndTime.Set(later),
	).Exec(ctx)
	handleError(err)

	// Add the raid bosses to a raid
	_, err = db.RaidBossesOnRaids.CreateOne(
		prisma.RaidBossesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID)),
		prisma.RaidBossesOnRaids.RaidBoss.Link(
			prisma.RaidBoss.ID.Equals(bossLichKing.ID)),
	).Exec(ctx)
	handleError(err)

	_, err = db.RaidBossesOnRaids.CreateOne(
		prisma.RaidBossesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid2.ID)),
		prisma.RaidBossesOnRaids.RaidBoss.Link(
			prisma.RaidBoss.ID.Equals(bossAlienQueen.ID)),
	).Exec(ctx)
	handleError(err)

	// Add the avatars to the first raid.
	_, err = db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Position.Set(1),
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar1.ID),
		),
	).Exec(ctx)
	handleError(err)

	_, err = db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Position.Set(2),
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar2.ID),
		),
	).Exec(ctx)
	handleError(err)

	// Create the events for the story
	// Days 1 through 6
	alienStoryEvent1 := seedEvent(db, alienStory.ID, 1)
	fmt.Println(alienStoryEvent1)

	alienStoryEvent2 := seedEvent(db, alienStory.ID, 2)
	fmt.Println(alienStoryEvent2)

	alienStoryEvent3 := seedEvent(db, alienStory.ID, 3)
	fmt.Println(alienStoryEvent3)

	alienStoryEvent4 := seedEvent(db, alienStory.ID, 4)
	fmt.Println(alienStoryEvent4)

	alienStoryEvent5 := seedEvent(db, alienStory.ID, 5)
	fmt.Println(alienStoryEvent5)

	alienStoryEvent6 := seedEvent(db, alienStory.ID, 6)
	fmt.Println(alienStoryEvent6)

	// Create the intro and completion messages.
	// Then the intro and completion messages to each event.

	// alien day 1 intro/completion messages
	seedEventMessage(db, alienStoryEvent1.ID, "The Sintari are a distant Alien hive-race that are distant relatives of Ants who have an almighty conquering Alien Queen with hoards of soldiers and minions serving her. They suck up resources and enslave species from every planet they encounter that serves their endless pre-programmed need for expansion. They now have their sights set on Earth, which has rare life-giving nutrients.\n\n"+
		"The CryptoFlowFightingForce (C3F) has spotted their scout ships at the edge of the galaxy and a forming a emergency team to defend earth from the impending invasion.  A call has gone out to the world's greatest Superheroes, scientists and soldiers to be the first wave of the resistance against an overwhelming threat...\n\n"+
		"Are you brave and courageous enough to join the C3F?\n",
		"The new heroes meet at Superhero HQ, don their new outfits, and receive training at the top secret facility at Area 52.\n\n",
	)

	// alien day 2 intro/completion messages
	seedEventMessage(db, alienStoryEvent2.ID, "Our heroes are honored to defend Earth and understand the sheer magnitude of this threat. The plan is to venture to the edge of the galaxy and run reconnaissance on the Scout ships gathering there. The far-reach spectral analysis has been unable to determine how many ships are forming on the edge of the galaxy and the size of the force threatening Earth.",
		"5... 4... 3... 2... 1. Blast off! The ship launches safely towards the moon base to refuel with our heroes onboard.",
	)

	// alien day 3 intro/completion messages
	seedEventMessage(db, alienStoryEvent3.ID, "Our heroes land safely on the moon using a SpaceZ 9000 rocket. They must refuel with the most powerful rocket fuel found only on the moon.",
		"Our heroes are debriefed at the super advanced, super secretive Lunar Base on the on the dark side the moon.",
	)

	// alien day 4 intro/completion messages
	seedEventMessage(db, alienStoryEvent4.ID, "Our heroes are further briefed by the Moon Galaxy Mission control on the details of their reconnaissance mission. It's a lot of information to take in as they need to cross navigate the galaxy, refueling at the rings of Saturday, and acquire some materials on Uranus before their final destination, the last sighting of the alien ships. This is going to be a long, quiet one.",
		"Our heroes pore their brains over the complexities required to cross navigate the length of the galaxy using hyper-flux technology.",
	)

	// alien day 5 intro/completion messages
	seedEventMessage(db, alienStoryEvent5.ID, "Emergency! A small battalion of unknown ships breaches the lunar orbit. The technology is unlike anything seen before. And these are presumably just the scout ships from the Alien Queen. Are our heroes ready?",
		"Our heroes are strapped in to defend. They weren't expecting to fight so soon. Is the team really ready for this?",
	)

	// alien day 6 intro/completion messages
	seedEventMessage(db, alienStoryEvent6.ID, "The aliens fire their advanced weaponry at the heroes' ship. They work as team to deftly maneuver among the onslaughts of beams! It's their turn now.",
		"Boom! The alien ships explodes in a silent purple-hued, ionized gas.  Our heroes have received their first taste of combat against the Sintari. This certainly won't be the last. They may not be so lucky next time. They realize they will need to learn to work together, encourage and motivate each other to work as strong team to save the Earth. They will need to learn many skills along the way and will need to be disciplined.",
	)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func toPtrString(str string) *string {
	return &str
}

func seedEvent(db *prisma.PrismaClient, storyID int, sequence int) prisma.EventModel {
	ctx := context.Background()

	evt, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(storyID)),
		prisma.Event.Sequence.Set(sequence),
	).Exec(ctx)
	handleError(err)

	return evt
}

func seedEventMessage(db *prisma.PrismaClient, eventId int, intro string, completion string) {
	ctx := context.Background()

	introMsg, err := db.Message.CreateOne(
		prisma.Message.Content.Set(intro),
		prisma.Message.Type.Set(messages.MessageTypeIntro.String()),
	).Exec(ctx)
	handleError(err)

	_, err = db.MessagesOnEvents.CreateOne(
		prisma.MessagesOnEvents.Message.Link(
			prisma.Message.ID.Equals(introMsg.ID)),
		prisma.MessagesOnEvents.Event.Link(
			prisma.Event.ID.Equals(eventId)),
	).Exec(ctx)

	completionMsg, err := db.Message.CreateOne(
		prisma.Message.Content.Set(completion),
		prisma.Message.Type.Set(messages.MessageTypeCompletion.String()),
	).Exec(ctx)
	handleError(err)

	_, err = db.MessagesOnEvents.CreateOne(
		prisma.MessagesOnEvents.Message.Link(
			prisma.Message.ID.Equals(completionMsg.ID)),
		prisma.MessagesOnEvents.Event.Link(
			prisma.Event.ID.Equals(eventId)),
	).Exec(ctx)
}
