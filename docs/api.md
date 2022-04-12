# API Reference

There are two APIs available and to participate in the game you need to implement both:

* [Websocket](#websocket)
* [REST](#rest)

## REST

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

You can use the REST API to list all available challenges

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

## Websocket
