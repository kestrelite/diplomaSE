# diplomaSE
A Diplomacy web app and adjudicator written in Go. It's free, I tell you, *free!*

## Wait, what?

#### What's Diplomacy?

[Diplomacy](https://en.wikibooks.org/wiki/Diplomacy/Rules) is a board game. It's most often described as Risk, except where turns happen simultaneously, and there's no luck involved. It requires alliance-making, and <s>sometimes</s> frequently involves backstabbing your friends. And who doesn't like *that*?

The rules are fairly simple. Your goal is to capture 18 supply centers by occupying them with your armies and fleets at the end of each game year. You can build units up to the number of supply centers you own. Armies can move to adjacent spaces on land; fleets can move to adjacent spaces on water. They can also support each other into occupied spaces, which gives the invading force more power against the resisting force. There are a handful more rules, and if you're curious, I suggest taking a look at [the more detailed rules](https://en.wikibooks.org/wiki/Diplomacy/Rules). 

#### What's a Diplomacy adjudicator?

Diplomacy rules, on the surface, are simple. Units can move from one square to an adjacent square, and units can support each other into those squares. In reality, though, algorithmically judging the rules is actually a rather finnicky task, and there are a significant number of errors first-release programs tend to make. In particular, there are some cases where establishing an order of precedence for the resolution of everone's orders simply doesn't work because there's a circular dependency. 

We'll skip over the details of what makes this tricky here, but if you'd like to get a sense for it, just scroll through the [Diplomacy Adjudicator Test Cases (DATC)](http://web.inter.nl.net/users/L.B.Kruijswijk/) or [The Math of Adjudication](http://www.diplom.org/Zine/S2009M/Kruijswijk/DipMath_Chp1.htm). (This program will be DATC-compliant.) 

#### What does this "free" thing mean?

**DiplomaSE** is intended to grow into a full client for Diplomacy games, and will ideally have a newfangled, easier-to-read Diplomacy game navigation system.

Since it's open-source, you're totally free to fork the repository and scroll through our code. Also, you'll be able to play on the full website once it's launched.

## Who created this thing?

This project grew out of a friendly/casual game of Diplomacy organized between
a small group of users on the Stack Exchange chat network, when one player
accidentally surrendered as a result of the silent failure of an AJAX request
that should have sent his orders to the server. That conversation, in a
nutshell:

> **A:** Maybe instead of playing the actual game, we should collaborate to build it
> a mobile app. That would be fun, and less anxiety-inducing.

> **U:** Or even our own webapp

> **E:** Actually, that could be a *ton* of fun.

> **E:** It'd be a shiny new way to try out Golang, too.

And thus, **diplomaSE** (working title) was born.

### Wait, so can I help?

Of course you can! That's why we're working on GitHub. That's why we wrote this
sweet README. *For you, our clamoring groupies.*

But seriously, contributions and input are welcome. You can find us in [the
SE Diplomacy chat room on Stack Exchange][1], or just send a pull request and
listen to us chortle at your noob code. Yes, chortle! Like immortal Wartortles
playing Portal. Fnordle. (**A** was drinking when he wrote this section.)

## When will it launch?

**DiplomaSE** is anticipated to launch in [6 to 8 weeks][2].

## Where <s>is the code located</s> am I?

**DiplomaSE** is hosted on GitHub. If you received this README in a shady deal
while trying to barter chickens for American health care or through some other
circumstance not foreseen by the authors, you can find the project home page
at https://github.com/telthien/diplomaSE.

## How could we build such a tool of evil!?

Apart from the motivations implied in the first section, we expect to gain the
following material and immaterial benefits as direct results of development:

- A fun introduction to, and/or better understanding of, the Go programming
language.
- Familiarity with some of the more peculiar cases of Diplomacy order
resolution.
- Rock-hard abs.
- [Making the world a better place.][3]

## Tools used

 - Go


  [1]: http://chat.stackexchange.com/rooms/27359/se-diplomacy
  [2]: http://meta.stackexchange.com/a/19514/254929
  [3]: https://vimeo.com/98720197
