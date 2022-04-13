# API Reference

There are two APIs available and to participate in the game you need to implement both:

* [REST](#rest)
  * [Registration](#registration)
  * [List available challenges](#list-available-challenges)
  * [Install a new challenge](#install-a-new-challenge)
* [Websocket](#websocket)
  * [join](#join)
  * [error](#error)
  * [attack](#attack)
  * [attack_challenge](#attack_challenge)
  * [attack_result](#attack_result)
  * [attack_solution](#attack_solution)
  * [defender_failed_to_defend](#defender_failed_to_defend)
  * [defend_action_request](#defend_action_request)
  * [defend_action](#defend_action)
  * [solution_evaluation_request](#solution_evaluation_request)
  * [solution_evaluation](#solution_evaluation)
* [Example usage](#example-usage)

## REST

You can find all DTOs [here](../core/engine/dto). Feel free to use this types while playing the game.

#### Registration

You can register yourself to one of the teams using this endpoint. The UUID you receive back should be persisted on your side throughout your game. If you loose your ID you also loose your progress.
NOTE: Your name has to be unique

```
POST /team/register
```

[Request body](../core/engine/dto/register.go):
```js
{
  name: "Unique name",
  team: "attacker" // or "defender"
}
```

[Response body](../core/engine/dto/register.go):
```js
{
  message: "Success",
  code: 200,
  id: "e256557a-e5c6-4475-a525-9857ea87cdad"
}
```

#### List available challenges

You can use this endpoint as an attacker to list all available challenges

```
GET /challenges
```

[Response body](../core/engine/dto/challenge.go):
```js
{
  message: "Success",
  code: 200,
  challenges: [
    {
      id: "fbb89d0f-3f11-43dc-a7fa-f31265df740b",
      name: "Reverse sorter",
      description: "You receive a random length string array in the first parameter of the hints. Your aim is to change the order of the array and send it back as the first parameter of the solution array",
      example: {
        hints: ["123456"],
        solutions: ["654321"]
      }
    }
  ]
}
```

#### Install a new challenge

You can use this endpoint as a defender to install new challenges in the system

NOTE: You will be expected to handle hint and solution evaluation requests as soon as you
installed a new module. So ideally you want to get ready for those steps, and this is the last step in your challenge creation flow.

```
POST /challenges
```

[Request body](../core/engine/dto/challenge.go)
```js
{
  name: "Reverse sorter",
  defenderId: "e256557a-e5c6-4475-a525-9857ea87cdad",
  description: "You receive a random length string array in the first parameter of the hints. Your aim is to change the order of the array and send it back as the first parameter of the solution array",
  example: {
    hints: ["123456"],
    solutions: ["654321"]
  }
}
```

[Response body](../core/engine/dto/challenge.go):
```js
{
  message: "Success",
  code: 200,
  id: "fbb89d0f-3f11-43dc-a7fa-f31265df740b"
}
```

You need to persist this ID from the response, as the system will use it to refer to your challenges when you are requested to provide hints or solution evaluations.

## Websocket

Websocket is available through `ws://host/team/join`. For this game we use [Gorilla Socket](https://github.com/gorilla/websocket) and it is recommended.

You can find all socket events typed [here](../core/engine/dto/socket.go).

You are encouraged to use this file either as a dependency or as a copy paste.

#### join

Emitted when the player is ready to join the live game.

Example message:
```js
{
  "type": "join",
  "id": "e256557a-e5c6-4475-a525-9857ea87cdad"
}
```

#### error

Emitted when an error happened during the flow.

Example message:
```js
{
  "type": "error",
  "message": "Could not parse Attack Event"
}
```

#### attack

Emitted to initiate an attack towards a challenge

Example message:
```js
{
  "type": "attack",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad"
}
```

#### attack_result

Emitted when an attacker receives results

Example message:
```js
{
  "type": "attack_result",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad",
  "success": true,
  "message": ""
}

{
  "type": "attack_result",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad",
  "success": false,
  "message": "Solutions array was too long"
}
```

#### attack_challenge

Emitted when hints are returned to the attacker

Example message:
```js
{
  "type": "attack_challenge",
  "hints": ["123456"],
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad"
}
```

#### attack_solution

Emitted when an attacker sends in solutions for hints

Example message:
```js
{
  "type": "attack_solution",
  "hints": ["123456"],
  "solutions": ["654321"],
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad"
}
```

#### defender_failed_to_defend

Emitted when a defender was offline either at requesting hints or at requesting solution validation

Example message:
```js
{
  "type": "defender_failed_to_defend",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad"
}
```

#### defend_action_request

Emitted when the defender is requested to provide hints for an attacker

Example message:
```js
{
  "type": "defend_action_request",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad",
  "combatId": "8049a606-6861-4536-8bcc-6449f50ae240"
}
```

#### defend_action

Emitted when the defender is providing hints for a challenge

Example message:
```js
{
  "type": "defend_action",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad",
  "combatId": "8049a606-6861-4536-8bcc-6449f50ae240",
  "hints": ["123456"]
}
```

#### solution_evaluation_request

Emitted when the defender is requested to evaluate a given solution

Example message:
```js
{
  "type": "solution_evaluation_request",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad",
  "combatId": "8049a606-6861-4536-8bcc-6449f50ae240",
  "hints": ["123456"],
  "solutions": ["654321"]
}
```

#### solution_evaluation

Emitted when the defender finished evaluation and ready to provide a result for the solutions

Example message:
```js
{
  "type": "solution_evaluation",
  "targetId": "e256557a-e5c6-4475-a525-9857ea87cdad",
  "combatId": "8049a606-6861-4536-8bcc-6449f50ae240",
  "success": false,
  "message": "Solutions array is too short"
}
```

## Example usage

You can find examples for attacking and defending [here](../example).
