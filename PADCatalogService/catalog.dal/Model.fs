namespace catalog.dal

open System
open catalog.Domain

[<AutoOpen>]
module Model =

  type Shoes =
    { Id: Guid
      Color: string
      Size: string
      Price: string
      Brand: string
      Category: string
      Model: string }

  module Shoes =
    open catalog

    let fromDomain (shoes: Domain.Shoes) =
      let (ShoesId id) = shoes.Id

      { Id = id
        Color = shoes.Color
        Size = shoes.Size
        Price = shoes.Price
        Brand = shoes.Brand
        Category = shoes.Category
        Model = shoes.Model }

    let toDomain (dto: Shoes) : Domain.Shoes =
      let shoesId = ShoesId dto.Id

      { Id = shoesId
        Color = dto.Color
        Size = dto.Size
        Price = dto.Price
        Brand = dto.Brand
        Category = dto.Category
        Model = dto.Model }
