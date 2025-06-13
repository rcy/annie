package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	"goirc/events"
	"goirc/handlers"
	"goirc/handlers/annie"
	"goirc/handlers/bedtime"
	"goirc/handlers/bible"
	"goirc/handlers/day"
	"goirc/handlers/ddate"
	"goirc/handlers/epigram"
	"goirc/handlers/gold"
	"goirc/handlers/hn"
	"goirc/handlers/kinfonet"
	"goirc/handlers/linkpool"
	"goirc/handlers/mlb"
	"goirc/handlers/tip"
	"goirc/handlers/tz"
	"goirc/handlers/weather"
	"goirc/internal/ai"
	"goirc/internal/cache"
	disco "goirc/internal/ddate"
	"goirc/internal/sun"
	db "goirc/model"
	"goirc/web"
	"log/slog"
	"regexp"
	"time"

	"github.com/robfig/cron"
)

func addHandlers(b *bot.Bot) {
	nick := regexp.QuoteMeta(b.Conn.GetNick())

	b.Handle(`^!catchup`, handlers.Catchup)
	b.Handle(`^,(.+)$`, handlers.CreateNote)
	b.Handle(`^([^\s:]+): (.+)$`, handlers.DeferredDelivery)
	b.Handle(`^!feedme`, handlers.AnonLink)
	b.Handle(`^!pipehealth\b`, handlers.AnonStatus)
	b.Handle(`(https?://\S+)`, handlers.Link)
	b.Handle(`^!day\b`, day.NationalDay)
	b.Handle(`^!dayi\b`, day.Dayi)
	b.Handle(`^!week\b`, day.NationalWeek)
	b.Handle(`^!weeki\b`, day.Weeki)
	b.Handle(`^!month\b`, day.NationalMonth)
	b.Handle(`^!monthi\b`, day.Monthi)
	b.Handle(`^!refs`, day.NationalRefs)
	b.Handle(`^!img (.+)$`, day.Image)
	b.Handle(`\b69[^0-9]*\b`, handlers.Nice)
	b.Handle(`^!odds`, mlb.PlayoffOdds)
	b.Handle(`^!godds`, mlb.GameOdds)
	b.Handle(`^!pom`, handlers.POM)
	b.Handle(`^("[^"]+)$`, handlers.Quote)
	b.Handle(`^!remindme ([^\s]+) (.+)$`, handlers.RemindMe)
	b.Handle(`^\?(\S+)`, handlers.Seen)
	b.Handle(`world.?cup`, handlers.Worldcup)
	b.Handle(`^!left`, handlers.TimeLeft)
	b.Handle(`^!epi`, epigram.Handle)
	b.Handle(`^!weather (.*)$`, weather.Handle)
	b.Handle(`^!weather$`, weather.Handle)
	b.Handle(`^!w (.*)$`, weather.Handle)
	b.Handle(`^!w$`, weather.Handle)
	b.Handle(`^!f (.*)$`, weather.HandleForecast)
	b.Handle(`^!f$`, weather.HandleForecast)
	b.Handle(`^!wf (.*)$`, weather.HandleWeatherForecast)
	b.Handle(`^!wf$`, weather.HandleWeatherForecast)
	b.Handle(`^!xweather (.+)$`, weather.XHandle)
	b.Handle(`^!k`, kinfonet.TodaysQuoteHandler)
	b.Handle(`^!gold`, gold.Handle)
	b.Handle(`^!hn`, hn.Handle)
	b.Handle(`^!auth$`, web.HandleAuth)
	b.Handle(`^!deauth$`, web.HandleDeauth)
	b.Handle(`night`, bedtime.Handle)
	b.Handle(fmt.Sprintf(`^%s:?(.+)$`, nick), annie.Handle)
	b.Handle(fmt.Sprintf(`^(.+),? %s.?$`, nick), annie.Handle)
	b.Handle(`^!bible (.+)$`, bible.Handle)
	b.Handle(`^tip$`, tip.Handle)
	b.Handle(`^date$`, ddate.Handle)
	b.Handle(`^tz`, tz.Handle)

	b.Repeat(10*time.Second, handlers.DoRemind)
	b.IdleRepeatAfterReset(8*time.Hour, handlers.POM)

	vancouver, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	toronto, err := time.LoadLocation("America/Toronto")
	if err != nil {
		panic(err)
	}

	q := model.New(db.DB.DB)

	c := cron.NewWithLocation(vancouver)
	err = c.AddFunc("16 14 15 * * 1,2,3,4,5,6", func() {
		note, err := q.RandomHistoricalTodayNote(context.TODO())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return
			}
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}

		b.Conn.Privmsgf(b.Channel, "on this day in %d, %s posted: %s", note.CreatedAt.Year(), note.Nick.String, note.Text.String)
	})
	if err != nil {
		panic(err)
	}

	// err = c.AddFunc("0 * * * * *", func() {
	// 	today := disco.NowIn(vancouver).WeekDay
	// 	if b.IsJoined {
	// 		if today == disco.SettingOrange {

	// 			b.Conn.Part(b.Channel)
	// 		}
	// 	} else {
	// 		if today != disco.SettingOrange {
	// 			b.Conn.Join(b.Channel)
	// 			time.Sleep(5 * time.Second)
	// 			if today == disco.Sweetmorn {
	// 				b.Conn.Privmsgf(b.Channel, "HAIL ERIS! GODDESS OF THE DAYS! LICK ME ON THIS SWEETMORN DAY! BE SURE I TASTE ALL NICE AND TASTY AND STUFF LIKE HOT FUDGE ON TOAST! SLURP!")
	// 			}
	// 		}
	// 	}
	// })
	// if err != nil {
	// 	log.Fatalf("c.AddFunc(discordia): %s", err)
	// }

	err = c.AddFunc("37 * * * * *", func() {
		ctx := context.TODO()
		zone := "America/Toronto"
		today := time.Now().In(toronto)
		todayValue := today.Format(time.DateOnly)

		rise, set, err := sun.SunriseSunset(today, zone, 43.64487, -79.38429)
		if err != nil {
			b.Conn.Privmsgf(b.Channel, "error: SunriseSunset: %s", err)
			return
		}
		slog.Debug("SunriseSunset", "rise", rise, "set", set, "tuset", time.Until(set))

		// any time in the last hour or in the next minute
		if time.Until(rise) >= -time.Hour && time.Until(rise) < time.Minute {
			value, err := cache.Get(ctx, "sunrise")
			if err != nil {
				b.Conn.Privmsgf(b.Channel, "error: cache.Get: %s", err)
				return
			}

			if value != todayValue {
				err := cache.Put(ctx, "sunrise", todayValue)
				if err != nil {
					b.Conn.Privmsgf(b.Channel, "error: cache.Put: %s", err)
				}
				events.Publish("sunrise", rise)
			}
		}

		// any time in the last hour or in the next minute
		if time.Until(set) >= -time.Hour && time.Until(set) < time.Minute {
			value, err := cache.Get(ctx, "sunset")
			if err != nil {
				b.Conn.Privmsgf(b.Channel, "error: cache.Get: %s", err)
				return
			}

			if value != todayValue {
				err := cache.Put(ctx, "sunset", todayValue)
				if err != nil {
					b.Conn.Privmsgf(b.Channel, "error: cache.Put: %s", err)
				}
				events.Publish("sunset", set)
			}
		}
	})

	err = c.AddFunc("57 * * * * *", func() {
		ctx := context.TODO()
		msg, err := q.ReadyFutureMessage(ctx, handlers.FutureMessageInterval)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return
			}
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}
		err = q.DeleteFutureMessage(ctx, msg.ID)
		if err != nil {
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}

		// send anonymous note
		switch msg.Kind {
		case "link":
			err = handlers.AnonLink(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
		case "quote":
			err = handlers.AnonQuote(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
		default:
			b.Conn.Privmsgf(b.Channel, "unhandled msg.Kind: %s", msg.Kind)
		}
		if err != nil {
			if errors.Is(err, ai.ErrBilling) {
				// the quote was sent, but no generated image, this is fine
				return
			}
			if errors.Is(err, linkpool.NoNoteFoundError) {
				// didn't find a note, reschedule
				_, scheduleErr := q.ScheduleFutureMessage(ctx, msg.Kind)
				if scheduleErr != nil {
					b.Conn.Privmsg(b.Channel, "error rescheduling: "+scheduleErr.Error())
				}
				return
			}
			// something else happened, spam the channel
			b.Conn.Privmsg(b.Channel, "error: "+err.Error())
		}
	})
	if err != nil {
		panic(err)
	}

	c.Start()

	events.Subscribe("anonnoteposted", func(note any) {
		go func() {
			err := handlers.AnonLink(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
			if err != nil {
				if errors.Is(err, ai.ErrBilling) {
					return
				}
				if errors.Is(err, linkpool.NoNoteFoundError) {
					return
				}
				b.Conn.Privmsg(b.Channel, "error: "+err.Error())
			}
		}()
	})

	events.Subscribe("anonquoteposted", func(note any) {
		go func() {
			err := handlers.AnonQuote(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
			if err != nil {
				if errors.Is(err, ai.ErrBilling) {
					return
				}
				if errors.Is(err, linkpool.NoNoteFoundError) {
					return
				}
				b.Conn.Privmsg(b.Channel, "error: "+err.Error())
			}
		}()
	})

	events.Subscribe("sunset", func(any) {
		if disco.FromTime(time.Now().In(toronto)).WeekDay == disco.SettingOrange {
			b.Conn.Privmsgf(b.Channel, "Hail Eris! Goddess of the Days! Look upon me as I look upon you on this Setting Orange day! I've had enough of this week, see you Sweetmorn! Good night!")
			b.Conn.Part(b.Channel)
		}
	})

	events.Subscribe("sunrise", func(any) {
		if disco.FromTime(time.Now().In(toronto)).WeekDay == disco.Sweetmorn {
			b.Conn.Join(b.Channel)
			time.Sleep(10 * time.Second)
			b.Conn.Privmsgf(b.Channel, "HAIL ERIS! GODDESS OF THE DAYS! LICK ME ON THIS SWEETMORN DAY! BE SURE I TASTE ALL NICE AND TASTY AND STUFF LIKE HOT FUDGE ON TOAST! SLURP!")
		}
	})

	b.Handle(`^!help`, func(params bot.HandlerParams) error {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, "https://github.com/rcy/annie/blob/main/handlers.go")
		return nil
	})
}
