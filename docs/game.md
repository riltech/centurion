# The game

You can find more detailed information about the flow and rules of the game here.

## Table of contents

* [Goal](#goal)
* [Game flow](#game-flow)
  * [Challenge](#challenge)
  * [Attackers flow](#attackers-flow)
  * [Defenders flow](#defenders-flow)
* [Point system](#point-system)
  * [Attackers](#attackers)
  * [Defenders](#defenders)

## Goal

The game consists of two teams, attackers and defenders. Defenders install challenges in Centurion and attackers are trying to solve these challenges. Each player are gathering individual points based on their activity, however, you also get points for reaching certain goals as a team.

## Game flow

Check the [API](./api.md) for more information about expected data structures.

The game starts by every player connecting to Centurion's API and registering a user. You will receive your ID after successful registration request, you **need to persist your ID throughout the game**. If you loose your ID, you loose your progress as well.

After registration, each team needs to connect to the live game through websocket using the previously acquired ID. 

#### Challenge

A challenge is essentially a function designed by defenders and reverse engineered by attackers.

There are two types of challenges:
* default - installed when the game starts
* player_created - installed by a given defender

Default modules are worth no point but they are a great way for attackers to design and test their code while supporting defenders with an example of how a challenge should be designed.

A challenge consists of the following information:

* Name - Name of the challenge
* Description - Detailed specification about how the solution should work
* Example - Containing an example hint and a solution for it

After a challenge is attacked, the defenders have to protect it by generating hints. This list can contain one or mulitple elements. For example, if the challenge is about changing element positions in the hints array, then it would contain multiple elements, but if the challenge is about data manipulation it might contain only one hint.

When hints are generated they are sent back to the attackers to provide solutions. This is also a list, that can contain one or multiple elements for the same reason as described above. After it is provided, the defender is requested to validate the solution. Individual points to players are given out after each outcome.

#### Attackers flow

As an attacker before you can solve challenges you need to acquire the currently available challenges using Centurion's rest API. Acquire the results then choose a challenge that you can solve.

Now the next step is to create your own code that in your best knowledge can solve the challenge. After you finished, you need to initiate an attack for the given challenge. You can do this using the websocket API. After this you will be either sent hints for the challenge by the defender or a result if the defender was unable to defend the module. (AKA was offline)

If you received hints, you need to solve it using the automation created before. Send it back for validation. The defender will be requested to validate, however, this will always end in a result for your direction.

Evaluate the results, if you failed, try to adjust your automation. If you succeeded congratulations, proceed to the next challenge!

#### Defenders flow

As a defender take a little time to observe the [default modules](../core/challenge/default.go) installed in Centurion. It will give you a good overview about how you should design challenges.

After you designed your first challenge, you need to install it using the REST API of Centurion.

Now a very **important difference** compared to an attacker is that you have to be able to defend your installed challenges. Which means, that after you installed your first challenge, you will be expected to stay online to provide hints and solution validations for attackers. Which means that the more resilient defender client you build the more individual points you can gain. You can only have one websocket session active at a time, so your goal is to reduce your downtime as much as possible.

When your challenge is being attacked, you will be requested to generate hint(s) for the challenge. This will be sent back to the attacker, then they need to provide solution(s) for the challenge. You need to validate the solution(s) provided for the challenge.

**The challenge you design has to be deterministic**. Which means that for a given input, we always need to get a given output. You solution validation will be tested by Centurion using your examples.


## Point system

You can find details here about how the point system works

#### Attackers

Individual scoring
* You get 1 point for every successful challenge resolution (you don't get points for invalid solutions) that you have not completed before
* You get 1 point for every 5 unique challenge solved
* You get 1 point for every attack where the defender has failed to defend

Team scoring

At the end of the game your team acquires extra points based on your ability to coordinate. A challenge is solved completely if everyone in your team was able to resolve it.
* You get 5 points if you (as a team) were able to solve at least 80% of the challenges
* You get 1 point for every challenge that were completed to a 100%

#### Defenders

Individual scoring
* You get 1 point if you have at least 1 challenge installed
* You get 1 point for every successful defense flow (even if the attack is successful)
* You get 1 point for every challenge solved, but not more than 50% of the attackers

Team scoring

At the end of the game your team acquires extra points based on your uptime as a team.
* \> 97% - 10 points
* \> 93% - 9 points
* \> 89% - 8 points
* \> 85% - 7 points
* \> 82% - 6 points
* \> 79% - 5 points
* \> 75% - 4 points
* \> 70% - 3 points
* \> 65% - 2 points
* < 65% - 1 point
