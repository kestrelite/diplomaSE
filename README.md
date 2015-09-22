# diplomaSE
A Diplomacy web app and adjudicator written in Go. It's free, I tell you, *free!*

## Who?

This project grew out of a friendly/casual game of Diplomacy organized between
a small group of users on the Stack Exchange chat network, when one player
accidentally surrendered as a result of the silent failure of an AJAX request
that should have sent his orders to the server (and silently failed instead). That conversation, in a
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

## What?

**DiplomaSE** is intended to grow into a full client for Diplomacy games. It is built with a custom DATC-compliant adjudicator written in Go, and will ideally have a newfangled, easier-to-read Diplomacy game navigation system. (Though it is not the same as [zond/godip](https://github.com/zond/godip), it relies on much of the same [math of adjudication](http://www.diplom.org/Zine/S2009M/Kruijswijk/DipMath_Chp1.htm).)

## When?

**DiplomaSE** is anticipated to launch in [6 to 8 weeks][2].

## Where?

**DiplomaSE** is hosted on GitHub. If you received this README in a shady deal
while trying to barter chickens for American health care or through some other
circumstance not foreseen by the authors, you can find the project home page
at https://github.com/telthien/diplomaSE.

## Why?

Apart from the motivations implied in the first section, we expect to gain the
following material and immaterial benefits as direct results of development:

- A fun introduction to, and/or better understanding of, the Go programming
language.
- Familiarity with some of the more peculiar cases of Diplomacy order
resolution.
- Rock-hard abs.
- [Making the world a better place.][3]

## How?

Three words: Go.


  [1]: http://chat.stackexchange.com/rooms/27359/se-diplomacy
  [2]: http://meta.stackexchange.com/a/19514/254929
  [3]: https://vimeo.com/98720197
