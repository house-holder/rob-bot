package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var triggers = []string{"reagan", "ronald", "nancy"}
var factTriggers = []string{"reagan facts", "reaganfacts", "reaganfax", "rgnfax", "fax", "fact", "facts"}
var facts = []string{
	"Reagan opposed the Civil Rights Act of 1964 and the Voting Rights Act of 1965.",
	"On August 5, 1981, Reagan fired 11,359 striking air traffic controllers and banned them from federal employment for life.",
	"Reagan's PATCO crackdown in 1981 is credited with emboldening corporate America to permanently replace striking workers.",
	"Reagan nearly tripled the national debt, from $998 billion in 1981 to $2.86 trillion in 1989.",
	"Reagan initially opposed making MLK Day a federal holiday, questioning whether King was 'important enough.'",
	"Reagan didn't publicly say the word \"AIDS\" until September 1985. By then, over 12,000 Americans had already died.",
	"Reagan cut funding for drug treatment programs while simultaneously massively expanding enforcement.",
	"Reagan raised taxes 11 times after his initial 1981 cut, including the largest peacetime tax increase in history in 1982.",
	"Reagan's deregulation contributed to the Savings & Loan crisis, which cost taxpayers $132 billion in bailouts.",
	"In 1971, Reagan called African UN diplomats 'monkeys' in a recorded call with Nixon. The tape was released in 2019.",
	"Reagan's first major speech on AIDS was May 31, 1987 - six years into the epidemic. Nearly 21,000 Americans were already dead.",
	"Reagan's War on Drugs increased the federal prison population from 24,000 in 1980 to over 100,000 by 1989.",
	"Reagan removed the solar panels Jimmy Carter installed on the White House roof in 1986.",
	"Reagan vetoed the Comprehensive Anti-Apartheid Act of 1986.",
	"After Reagan fired PATCO air traffic controllers in 1981, it took nearly a decade to restore staffing levels.",
	"Reagan's 1980 campaign launched in Philadelphia, MS (where three civil rights workers were murdered in 1964) with mentions of \"states' rights.\"",
	"Reagan closed psychiatric hospitals as California governor in 1967, contributing directly to modern homelessness in the state.",
	"When asked about AIDS in 1982, Reagan's press secretary Larry Speakes laughed and called it 'the gay plague.'",
	"Reagan's administration trained and funded the Salvadoran military, which massacred over 800 civilians at El Mozote in 1981.",
	"The 1986 Anti-Drug Abuse Act, signed by Reagan, created the 100:1 crack vs. powder cocaine sentencing disparity—targeting Black communities.",
	"Reagan's EPA budget was cut by 22%. Staff was reduced by 20%.",
	"In 1981, Reagan claimed trees cause more pollution than automobiles.",
	"Reagan's Interior Secretary James Watt said he didn't worry about protecting the environment because Jesus would return soon.",
	"In 1961, Reagan recorded an LP warning that Medicare would lead to socialism and the end of freedom in America.",
	"Reagan eliminated the FCC Fairness Doctrine in 1987, enabling the rise of one-sided partisan talk radio.",
	"Reagan's HUD Secretary Samuel Pierce oversaw a scandal where housing funds were steered to Republican donors. 16 officials were convicted.",
	"In 1981, Reagan's administration tried to classify ketchup as a vegetable to cut school lunch costs.",
	"Reagan slashed federal housing assistance by 75%, contributing directly to the 1980s homelessness crisis.",
	"The poverty rate rose from 11.7% to 15.2% during Reagan's first term.",
	"Reagan cut the top marginal tax rate from 70% to 28%, the largest tax cut for the wealthy in American history.",
	"Real wages for the bottom 50% of workers declined during the Reagan years despite GDP growth.",
	"Reagan didn't just fire the PATCO strikers — he decertified the entire union and banned the workers from federal jobs for life.",
	"PATCO had endorsed Reagan in the 1980 election. He crushed them anyway.",
	"After Reagan's PATCO bust, the use of permanent replacement workers by private employers increased dramatically throughout the 1980s.",
	"Reagan stacked the National Labor Relations Board with anti-union appointees. Under his NLRB, enforcement of employer violations dropped sharply.",
	"Reagan's NLRB Chairman Donald Dotson openly stated that collective bargaining \"destroys individual freedom.\"",
	"Before Reagan, it was rare for employers to permanently replace strikers. After PATCO, it became standard corporate practice.",
	"Under Reagan's NLRB, complaints settled in favor of laborers dropped by more than 50%.",
	"Reagan closed one-third of OSHA field offices and cut staff by over 25%.",
	"Under Reagan, OSHA penalties against employers dropped by nearly 75%.",
	"Reagan's OSHA shifted from enforcement to \"voluntary compliance\", preferring that employers police themselves.",
	"Union membership fell from 23% to 16% of the workforce during Reagan's presidency.",
	"Reagan made efforts to weaken child labor laws and create a sub-minimum wage for young workers.",
	"Reagan's 1982 Bus Regulatory Reform Act led directly to the 1983 Greyhound strike, which ultimately decimated bus driver unions.",
}

func msgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // detect self-messages and ignore
		return
	}
	msg := strings.ToLower(m.Content)

	for _, ft := range factTriggers {
		if strings.Contains(msg, ft) {
			fact := facts[rand.Intn(len(facts))]
			fullMessage := fmt.Sprintf(">>> %s", fact)
			s.ChannelMessageSend(m.ChannelID, fullMessage)
			return
		}
	}

	for _, t := range triggers {
		if strings.Contains(msg, t) {
			if t == "nancy" {
				s.ChannelMessageSend(m.ChannelID, "fuck the reagans")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "fuck ronald reagan")
			return
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env load fail: ", err)
	}

	token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Create Discord session fail: ", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	dg.AddHandler(msgCreate)
	dg.Open()
	log.Println("Bot is running. q to quit, r to reload.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}
