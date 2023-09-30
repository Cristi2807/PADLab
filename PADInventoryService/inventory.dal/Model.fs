namespace inventory.dal

open System
open inventory.Domain

[<AutoOpen>]
module Model =

  type Transaction =
    { Id: Guid
      ShoesId: Guid
      CreationDate: DateTime
      Quantity: int64
      OperationType: int16 }

  module Transaction =
    open inventory

    let fromDomain (transaction: Domain.Transaction) =
      let (TransactionId id) = transaction.Id
      let (ShoesId shoesId) = transaction.ShoesId

      { Id = id
        ShoesId = shoesId
        CreationDate = transaction.CreationDate
        Quantity = transaction.Quantity
        OperationType = transaction.OperationType }

    let toDomain (dto: Transaction) : Domain.Transaction =
      let id = TransactionId dto.Id
      let shoesId = ShoesId dto.ShoesId

      { Id = id
        ShoesId = shoesId
        CreationDate = dto.CreationDate
        Quantity = dto.Quantity
        OperationType = dto.OperationType }
