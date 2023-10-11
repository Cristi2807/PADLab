namespace catalog.dal

open Npgsql
open System
open FsToolkit.ErrorHandling
open catalog.ServiceUtils

type DalContext =
  { DbConnection: NpgsqlConnection
    DbTransaction: NpgsqlTransaction option }

  interface IDisposable with
    member this.Dispose() =
      try
        match this.DbTransaction with
        | Some tr -> tr.Dispose()
        | _ -> ()

        this.DbConnection.Dispose()
      with _ ->
        ()

type DalContextFactory() =
  let buildConnString =
    sprintf "Server=%s; Port=%d; Database=%s; User Id=%s; Password=%s;"

  let mutable writeCtx: DalContext option = None

  member private __.create connString =
    let conn = new NpgsqlConnection(connString)

    conn.OpenAsync()
    |> UnitTaskToAsyncResult
    |> AsyncResult.map (fun () ->
      { DbConnection = conn
        DbTransaction = None })

  member private __.createWithTransaction connString =
    __.create connString
    |> AsyncResult.map (fun ctx ->
      let tr = ctx.DbConnection.BeginTransaction()
      { ctx with DbTransaction = Some tr })

  member private __.getEnvVariable varName =
    asyncResult {
      match Environment.GetEnvironmentVariable varName with
      | null ->
        return!
          $"{varName} Environment Variable not set!"
          |> InternalError
          |> AsyncResult.returnError
      | x -> return x
    }

  member __.CommitTransaction() =
    match writeCtx with
    | Some dalCtx ->
      match dalCtx.DbTransaction with
      | Some tr ->
        tr.Commit()
        tr.Dispose()
        (dalCtx :> IDisposable).Dispose()
      | _ -> ()
    | _ -> ()

  member __.Rollback() =
    match writeCtx with
    | Some dalCtx ->
      match dalCtx.DbTransaction with
      | Some tr ->
        tr.Rollback()
        tr.Dispose()
        (dalCtx :> IDisposable).Dispose()
      | _ -> ()
    | _ -> ()

  member __.GetReadContext() =
    asyncResult {
      let! dbHost = __.getEnvVariable "DB_RO_HOST"
      let! dbName = __.getEnvVariable "DB_NAME"
      let! dbUser = __.getEnvVariable "DB_USER"
      let! dbPassword = __.getEnvVariable "DB_PASSWORD"
      let connString = buildConnString dbHost 5432 dbName dbUser dbPassword
      return! __.create connString
    }

  member __.GetWriteContext() =
    asyncResult {
      match writeCtx with
      | Some ctx -> return ctx
      | _ ->
        let! dbHost = __.getEnvVariable "DB_HOST"
        let! dbName = __.getEnvVariable "DB_NAME"
        let! dbUser = __.getEnvVariable "DB_USER"
        let! dbPassword = __.getEnvVariable "DB_PASSWORD"
        let connString = buildConnString dbHost 5432 dbName dbUser dbPassword
        let! newCtx = __.createWithTransaction connString
        writeCtx <- Some newCtx
        return newCtx
    }

  member __.HandleTransaction(res: Async<Result<string, ServiceError>>) =
    res
    |> AsyncResult.tee (fun _ ->
      __.CommitTransaction()
      __.Dispose())
    |> AsyncResult.teeError (fun _ ->
      __.Rollback()
      __.Dispose())

  member __.Dispose() =
    writeCtx |> Option.iter (fun ctx -> (ctx :> IDisposable).Dispose())

  interface IDisposable with
    member this.Dispose() = this.Dispose()
