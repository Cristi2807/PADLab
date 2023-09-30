namespace catalog

open System

module Domain =

  type ShoesId = ShoesId of Guid

  type Shoes =
    { Id: ShoesId
      Color: string
      Size: string
      Price: string
      Brand: string
      Category: string
      Model: string }


open Domain

module Decode =
  open Thoth.Json.Net

  let shoes () : Decoder<Shoes> =
    Decode.object (fun get ->
      { Shoes.Id = get.Optional.Field "id" Decode.guid |> Option.defaultValue Guid.Empty |> ShoesId
        Shoes.Color = get.Required.Field "color" Decode.string
        Shoes.Size = get.Required.Field "size" Decode.string
        Shoes.Price = get.Required.Field "price" Decode.string
        Shoes.Brand = get.Required.Field "brand" Decode.string
        Shoes.Category = get.Required.Field "category" Decode.string
        Shoes.Model = get.Required.Field "model" Decode.string })

module Encode =
  open Thoth.Json.Net

  let ShoesId (ShoesId shoesId) = Encode.guid shoesId

  let shoes (shoes: Shoes) =
    [ "id", ShoesId shoes.Id
      "color", Encode.string shoes.Color
      "size", Encode.string shoes.Size
      "price", Encode.string shoes.Price
      "brand", Encode.string shoes.Brand
      "category", Encode.string shoes.Category
      "model", Encode.string shoes.Model ]
    |> Encode.object
